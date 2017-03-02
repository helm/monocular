package helpers

import (
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/helm/monocular/src/api/config"
	"github.com/helm/monocular/src/api/swagger/models"
	"gopkg.in/yaml.v2"

	"github.com/helm/monocular/src/api/data/cache/charthelper"
)

// APIVer1String is the API version 1 string we include in route URLs
const APIVer1String = "v1"

// IsYAML checks for valid YAML
func IsYAML(b []byte) bool {
	var yml map[string]interface{}
	ret := yaml.Unmarshal(b, &yml)
	return ret == nil
}

// ParseYAMLRepo converts a YAML representation of a repo
// to a slice of charts
func ParseYAMLRepo(rawYAML []byte, repoName string) ([]*models.ChartPackage, error) {
	var ret []*models.ChartPackage
	repoIndex := make(map[interface{}]interface{})
	if err := yaml.Unmarshal(rawYAML, &repoIndex); err != nil {
		return nil, err
	}
	entries := repoIndex["entries"]
	if entries == nil {
		return nil, fmt.Errorf("error parsing entries from YAMLified repo")
	}
	e, _ := yaml.Marshal(&entries)
	chartEntries := make(map[string][]models.ChartPackage)
	if err := yaml.Unmarshal(e, &chartEntries); err != nil {
		return nil, err
	}
	for entry := range chartEntries {
		for i := range chartEntries[entry] {
			chartEntries[entry][i].Repo = repoName
			ret = append(ret, &chartEntries[entry][i])
		}
	}
	return ret, nil
}

// MakeChartResource composes a Resource type that represents a repo+chart
func MakeChartResource(chart *models.ChartPackage) *models.Resource {
	var ret models.Resource
	ret.Type = StrToPtr("chart")
	ret.ID = StrToPtr(MakeChartID(chart.Repo, *chart.Name))
	ret.Attributes = &models.Chart{
		Repo:        getRepoObject(chart),
		Name:        chart.Name,
		Description: chart.Description,
		Home:        chart.Home,
		Sources:     chart.Sources,
		Keywords:    chart.Keywords,
		Maintainers: chart.Maintainers,
	}
	AddLatestChartVersionRelationship(&ret, chart)
	return &ret
}

// MakeChartResources accepts a slice of repo+chart data, converts each to a Resource type
// and then returns the slice of the converted Resource types (throwing away version information,
// and collapsing all chart+version records into a single resource representation for each chart)
func MakeChartResources(charts []*models.ChartPackage) []*models.Resource {
	var chartsResource []*models.Resource
	found := make(map[string]bool)
	for _, chart := range charts {
		if !found[*chart.Name] {
			found[*chart.Name] = true
			latestVersion, _ := GetLatestChartVersion(charts, *chart.Name)
			resource := MakeChartResource(latestVersion)
			AddCanonicalLink(resource)
			AddLatestChartVersionRelationship(resource, latestVersion)
			chartsResource = append(chartsResource, resource)
		}
	}
	return chartsResource
}

// MakeChartVersionResource composes a Resource type that represents a chartVersion
func MakeChartVersionResource(chart *models.ChartPackage) *models.Resource {
	var ret models.Resource
	ret.Type = StrToPtr("chartVersion")
	ret.ID = StrToPtr(MakeChartVersionID(chart.Repo, *chart.Name, *chart.Version))
	ret.Attributes = &models.ChartVersion{
		Created: chart.Created,
		Digest:  chart.Digest,
		Urls:    chart.Urls,
		Version: chart.Version,
		Icons:   makeAvailableIcons(chart),
		Readme:  makeReadmeURL(chart),
	}
	AddChartRelationship(&ret, chart)
	return &ret
}

// MakeChartVersionResources accepts a slice of versioned repo+chart data, converts each to a Resource type
// and then returns the slice of the converted Resource types (retaining version info)
func MakeChartVersionResources(charts []*models.ChartPackage) []*models.Resource {
	var chartsResource []*models.Resource
	for _, chart := range charts {
		resource := MakeChartVersionResource(chart)
		chartsResource = append(chartsResource, resource)
	}
	return chartsResource
}

// AddChartRelationship adds a "relationships" reference to a chartVersion resource's chart
func AddChartRelationship(resource *models.Resource, chartPackage *models.ChartPackage) {
	resource.Relationships = &models.ChartRelationship{
		Chart: &models.ChartAsRelationship{
			Links: &models.ResourceLink{
				Self: StrToPtr(MakeRepoChartRouteURL(APIVer1String, chartPackage.Repo, *chartPackage.Name)),
			},
			Data: &models.Chart{
				Name:        chartPackage.Name,
				Description: chartPackage.Description,
				Repo:        getRepoObject(chartPackage),
				Home:        chartPackage.Home,
				Sources:     chartPackage.Sources,
				Maintainers: chartPackage.Maintainers,
			},
		},
	}
}

// AddLatestChartVersionRelationship adds a "relationships" reference to a chart resource's latest chartVersion
func AddLatestChartVersionRelationship(resource *models.Resource, chartPackage *models.ChartPackage) {
	resource.Relationships = &models.LatestChartVersionRelationship{
		LatestChartVersion: &models.ChartVersionAsRelationship{
			Links: &models.ResourceLink{
				Self: StrToPtr(MakeRepoChartVersionRouteURL(APIVer1String, chartPackage.Repo, *chartPackage.Name, *chartPackage.Version)),
			},
			Data: &models.ChartVersion{
				Created: chartPackage.Created,
				Digest:  chartPackage.Digest,
				Urls:    chartPackage.Urls,
				Version: chartPackage.Version,
				Icons:   makeAvailableIcons(chartPackage),
				Readme:  makeReadmeURL(chartPackage),
			},
		},
	}
}

// AddCanonicalLink adds a "self" link to a chart resource's canonical API endpoint
func AddCanonicalLink(resource *models.Resource) {
	resource.Links = &models.ResourceLink{
		Self: StrToPtr(MakeRepoChartRouteURL(APIVer1String, *resource.Attributes.(*models.Chart).Repo.Name, *resource.Attributes.(*models.Chart).Name)),
	}
}

// GetLatestChartVersion returns the most recent version from a slice of versioned charts
func GetLatestChartVersion(charts []*models.ChartPackage, name string) (*models.ChartPackage, error) {
	latest := "0.0.0"
	var ret *models.ChartPackage
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
func GetChartVersion(charts []*models.ChartPackage, name, version string) (*models.ChartPackage, error) {
	var ret *models.ChartPackage
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
func GetChartVersions(charts []*models.ChartPackage, name string) ([]*models.ChartPackage, error) {
	var ret []*models.ChartPackage
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

// MakeRepoChartRouteURL returns a string that represents
// /{:apiVersion}/charts/{:repo}/{:chart}
func MakeRepoChartRouteURL(apiVer, repo, name string) string {
	return fmt.Sprintf("/%s/charts/%s/%s", apiVer, repo, name)
}

// MakeRepoChartVersionRouteURL returns a string that represents
// /{:apiVersion}/charts/{:repo}/{:chart}/versions/{:version}
func MakeRepoChartVersionRouteURL(apiVer, repo, name, version string) string {
	return fmt.Sprintf("/%s/charts/%s/%s/versions/%s", apiVer, repo, name, version)
}

// MakeChartID returns a chart ID in the form {:repo}/{:chart}
func MakeChartID(repo, chart string) string {
	return fmt.Sprintf("%s/%s", repo, chart)
}

// MakeChartVersionID returns a chartVersion ID in the form {:repo}/{:chart}:{:version}
func MakeChartVersionID(repo, chart, version string) string {
	return fmt.Sprintf("%s/%s:%s", repo, chart, version)
}

// newestSemVer returns the newest (largest) semver string
func newestSemVer(v1 string, v2 string) (string, error) {
	v1Semver, err := semver.NewVersion(v1)
	if err != nil {
		return "", err
	}

	v2Semver, err := semver.NewVersion(v2)
	if err != nil {
		return "", err
	}

	if v1Semver.LessThan(v2Semver) {
		return v2, nil
	}
	return v1, nil
}

// Int64ToPtr converts an int64 to an *int64
func Int64ToPtr(n int64) *int64 {
	return &n
}

// StrToPtr converts a string to a *string
func StrToPtr(s string) *string {
	return &s
}

func makeAvailableIcons(chart *models.ChartPackage) []*models.Icon {
	var res []*models.Icon
	icons := charthelper.AvailableIcons(chart, "/assets")
	for _, icon := range icons {
		res = append(res, &models.Icon{Name: &icon.Name, Path: &icon.Path})
	}
	return res
}

func makeReadmeURL(chart *models.ChartPackage) *string {
	res := charthelper.ReadmeStaticUrl(chart, "/assets")
	return &res
}

func getRepoObject(chart *models.ChartPackage) *models.Repo {
	var repoPayload models.Repo

	config, _ := config.GetConfig()
	for _, repo := range config.Repos {
		if repo.Name == chart.Repo {
			repoPayload = models.Repo{
				Name:   &repo.Name,
				URL:    &repo.URL,
				Source: repo.Source,
			}
			return &repoPayload
		}
	}
	return &repoPayload
}
