package helpers

import (
	"fmt"
	"strings"

	"github.com/helm/monocular/src/api/swagger/models"
	"gopkg.in/yaml.v2"
)

const apiVer1 = "v1"

// IsYAML checks for valid YAML
func IsYAML(b []byte) bool {
	var yml map[string]interface{}
	ret := yaml.Unmarshal(b, &yml)
	return ret == nil
}

// ParseYAMLRepo converts a YAML representation of a repo
// to a slice of versioned charts
func ParseYAMLRepo(rawYAML []byte) ([]*models.ChartVersion, error) {
	repo := make(map[interface{}]interface{})
	if err := yaml.Unmarshal(rawYAML, &repo); err != nil {
		return nil, err
	}
	entries := repo["entries"]
	if entries == nil {
		return nil, fmt.Errorf("error parsing entries from YAMLified repo")
	}
	e, _ := yaml.Marshal(&entries)
	chartEntries := make(map[string][]models.ChartVersion)
	if err := yaml.Unmarshal(e, &chartEntries); err != nil {
		return nil, err
	}
	var charts []*models.ChartVersion
	for entry := range chartEntries {
		for i := range chartEntries[entry] {
			charts = append(charts, &chartEntries[entry][i])
		}
	}
	return charts, nil
}

// MakeChartResource composes a Resource type that represents a repo+chart
func MakeChartResource(chart *models.ChartVersion, repo string) *models.Resource {
	var ret models.Resource
	ret.Type = StrToPtr("chart")
	ret.ID = StrToPtr(fmt.Sprintf("%s/%s", repo, *chart.Name))
	ret.Attributes = &models.ChartResourceAttributes{
		Repo:        &repo,
		Name:        chart.Name,
		Description: chart.Description,
		Home:        chart.Home,
		Sources:     chart.Sources,
		Keywords:    chart.Keywords,
		Maintainers: chart.Maintainers,
		Icon:        chart.Icon,
	}
	return &ret
}

// MakeChartsResource accepts a slice of versioned repo+chart data, converts each to a Resource type
// and then returns the slice of the converted Resource types (throwing away version information,
// and collapsing all chart+version records into a single resource representation for each chart)
func MakeChartsResource(charts []*models.ChartVersion, repo string) []*models.Resource {
	var chartsResource []*models.Resource
	found := make(map[string]bool)
	for _, chart := range charts {
		if !found[*chart.Name] {
			found[*chart.Name] = true
			resource := MakeChartResource(chart, repo)
			AddCanonicalLink(resource)
			chartsResource = append(chartsResource, resource)
		}
	}
	return chartsResource
}

// MakeChartVersionResource composes a Resource type that represents a repo+chart
func MakeChartVersionResource(chart *models.ChartVersion, repo string) *models.Resource {
	var ret models.Resource
	ret.Type = StrToPtr("chartVersion")
	ret.ID = StrToPtr(fmt.Sprintf("%s/%s:%s", repo, *chart.Name, *chart.Version))
	ret.Attributes = &models.ChartVersionResourceAttributes{
		Repo:        &repo,
		Name:        chart.Name,
		Version:     chart.Version,
		Description: chart.Description,
		Created:     chart.Created,
		Digest:      chart.Digest,
		Home:        chart.Home,
		Sources:     chart.Sources,
		Urls:        chart.Urls,
		Keywords:    chart.Keywords,
		Maintainers: chart.Maintainers,
		Icon:        chart.Icon,
	}
	return &ret
}

// MakeChartVersionsResource accepts a slice of versioned repo+chart data, converts each to a Resource type
// and then returns the slice of the converted Resource types (retaining version info)
func MakeChartVersionsResource(charts []*models.ChartVersion, repo string) []*models.Resource {
	var chartsResource []*models.Resource
	for _, chart := range charts {
		resource := MakeChartVersionResource(chart, repo)
		chartsResource = append(chartsResource, resource)
	}
	return chartsResource
}

// AddLatestRelationship adds a "relationships" reference to the resource's latest chartVersion release
func AddLatestRelationship(resource *models.Resource, chartVersion *models.ChartVersion) {
	resource.Relationships = &models.ChartLatestRelationships{
		ChartVersionLatest: &models.ChartVersionRelationshipAttributes{
			Links: &models.ChartLatestLinks{
				Related: StrToPtr(fmt.Sprintf("/%s/charts/%s/%s/versions/%s", apiVer1, *resource.Attributes.(*models.ChartResourceAttributes).Repo, *resource.Attributes.(*models.ChartResourceAttributes).Name, *chartVersion.Version)),
			},
			Data: &models.ChartVersionLatestResourceAttributes{
				Created: chartVersion.Created,
				Digest:  chartVersion.Digest,
				Urls:    chartVersion.Urls,
				Version: chartVersion.Version,
			},
		},
	}
}

// AddCanonicalLink adds a "canonical" reference to the resource's "links" object
func AddCanonicalLink(resource *models.Resource) {
	resource.Links = &models.ChartResourceLinks{
		Canonical: StrToPtr(fmt.Sprintf("/%s/charts/%s/%s", apiVer1, *resource.Attributes.(*models.ChartResourceAttributes).Repo, *resource.Attributes.(*models.ChartResourceAttributes).Name)),
	}
}

// GetLatestChartVersion returns the most recent version from a slice of versioned charts
func GetLatestChartVersion(charts []*models.ChartVersion, name string) (*models.ChartVersion, error) {
	latest := "0.0.0"
	var ret *models.ChartVersion
	for _, chart := range charts {
		if *chart.Name == name {
			newest, err := newestSemVer(latest, *chart.Version)
			if err != nil {
				return nil, err
			}
			latest = newest
			if latest == *chart.Version {
				ret = chart
			}
		}
	}
	if ret == nil {
		return ret, fmt.Errorf("chart %s not found\n", name)
	}
	return ret, nil
}

// GetChartVersion returns a specific versions of a chart
func GetChartVersion(charts []*models.ChartVersion, name, version string) (*models.ChartVersion, error) {
	var ret *models.ChartVersion
	for _, chart := range charts {
		if *chart.Name == name && *chart.Version == version {
			ret = chart
		}
	}
	if ret == nil {
		return ret, fmt.Errorf("didn't find version %s of chart %s\n", version, name)
	}
	return ret, nil
}

// GetChartVersions returns all versions of a chart
func GetChartVersions(charts []*models.ChartVersion, name string) ([]*models.ChartVersion, error) {
	var ret []*models.ChartVersion
	for _, chart := range charts {
		if *chart.Name == name {
			ret = append(ret, chart)
		}
	}
	if ret == nil {
		return ret, fmt.Errorf("no chart versions found for %s\n", name)
	}
	return ret, nil
}

// newestSemVer returns the newest (largest) semver string
func newestSemVer(v1 string, v2 string) (string, error) {
	v1Slice := strings.Split(v1, ".")
	if len(v1Slice) != 3 {
		return "", semverStringError(v1)
	}
	v2Slice := strings.Split(v2, ".")
	if len(v2Slice) != 3 {
		return "", semverStringError(v2)
	}
	for i, subVer1 := range v1Slice {
		if v2Slice[i] > subVer1 {
			return v2, nil
		} else if subVer1 > v2Slice[i] {
			return v1, nil
		}
	}
	return v1, nil
}

// semverStringError returns a bad semver string error
func semverStringError(v string) error {
	return fmt.Errorf("%s is not a semver-compatible string", v)
}

// Int64ToPtr converts an int64 to an *int64
func Int64ToPtr(n int64) *int64 {
	return &n
}

// StrToPtr converts a string to a *string
func StrToPtr(s string) *string {
	return &s
}
