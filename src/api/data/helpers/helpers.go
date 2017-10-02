package helpers

import (
	"fmt"

	"github.com/Masterminds/semver"
	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
	"github.com/kubernetes-helm/monocular/src/api/datastore"
	"github.com/kubernetes-helm/monocular/src/api/models"
	swaggermodels "github.com/kubernetes-helm/monocular/src/api/swagger/models"

	"github.com/kubernetes-helm/monocular/src/api/data/cache/charthelper"
	"github.com/kubernetes-helm/monocular/src/api/data/pointerto"
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
func ParseYAMLRepo(rawYAML []byte, repoName string) ([]*swaggermodels.ChartPackage, error) {
	var ret []*swaggermodels.ChartPackage
	repoIndex := make(map[string]interface{})
	if err := yaml.Unmarshal(rawYAML, &repoIndex); err != nil {
		return nil, err
	}
	entries := repoIndex["entries"]
	if entries == nil {
		return nil, fmt.Errorf("error parsing entries from YAMLified repo")
	}
	e, _ := yaml.Marshal(&entries)
	chartEntries := make(map[string][]swaggermodels.ChartPackage)
	if err := yaml.Unmarshal(e, &chartEntries); err != nil {
		return nil, err
	}
	for entry := range chartEntries {
		if chartEntries[entry][0].Deprecated != nil && *chartEntries[entry][0].Deprecated {
			log.WithFields(log.Fields{
				"name": entry,
			}).Info("Deprecated chart skipped")
			continue
		}
		for i := range chartEntries[entry] {
			chartEntries[entry][i].Repo = repoName
			ret = append(ret, &chartEntries[entry][i])
		}
	}
	return ret, nil
}

// MakeChartResource composes a Resource type that represents a repo+chart
func MakeChartResource(db datastore.Database, chart *swaggermodels.ChartPackage) *swaggermodels.Resource {
	var ret swaggermodels.Resource
	ret.Type = pointerto.String("chart")
	ret.ID = pointerto.String(MakeChartID(chart.Repo, *chart.Name))
	ret.Attributes = &swaggermodels.Chart{
		Repo:        getRepoObject(db, chart),
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

// MakeRepoResource composes a Resource type that represents a repository
func MakeRepoResource(repo *models.Repo) *swaggermodels.Resource {
	var ret swaggermodels.Resource
	ret.Type = pointerto.String("repository")
	ret.ID = &repo.Name
	ret.Attributes = repo
	return &ret
}

// MakeRepoResources returns an array of RepoResources
func MakeRepoResources(repos []*models.Repo) []*swaggermodels.Resource {
	var reposResource []*swaggermodels.Resource
	for _, repo := range repos {
		resource := MakeRepoResource(repo)
		reposResource = append(reposResource, resource)
	}
	return reposResource
}

// MakeChartResources accepts a slice of repo+chart data, converts each to a Resource type
// and then returns the slice of the converted Resource types (throwing away version information,
// and collapsing all chart+version records into a single resource representation for each chart)
func MakeChartResources(db datastore.Database, charts []*swaggermodels.ChartPackage) []*swaggermodels.Resource {
	var chartsResource []*swaggermodels.Resource
	found := make(map[string]bool)
	for _, chart := range charts {
		if !found[*chart.Name] {
			found[*chart.Name] = true
			latestVersion, _ := GetLatestChartVersion(charts, *chart.Name)
			resource := MakeChartResource(db, latestVersion)
			AddCanonicalLink(resource)
			AddLatestChartVersionRelationship(resource, latestVersion)
			chartsResource = append(chartsResource, resource)
		}
	}
	return chartsResource
}

// MakeChartVersionResource composes a Resource type that represents a chartVersion
func MakeChartVersionResource(db datastore.Database, chart *swaggermodels.ChartPackage) *swaggermodels.Resource {
	var ret swaggermodels.Resource
	ret.Type = pointerto.String("chartVersion")
	ret.ID = pointerto.String(MakeChartVersionID(chart.Repo, *chart.Name, *chart.Version))
	ret.Attributes = &swaggermodels.ChartVersion{
		Created:    chart.Created,
		Digest:     chart.Digest,
		Urls:       chart.Urls,
		Version:    chart.Version,
		AppVersion: chart.AppVersion,
		Icons:      makeAvailableIcons(chart),
		Readme:     makeReadmeURL(chart),
	}
	AddChartRelationship(db, &ret, chart)
	return &ret
}

// MakeChartVersionResources accepts a slice of versioned repo+chart data, converts each to a Resource type
// and then returns the slice of the converted Resource types (retaining version info)
func MakeChartVersionResources(db datastore.Database, charts []*swaggermodels.ChartPackage) []*swaggermodels.Resource {
	var chartsResource []*swaggermodels.Resource
	for _, chart := range charts {
		resource := MakeChartVersionResource(db, chart)
		chartsResource = append(chartsResource, resource)
	}
	return chartsResource
}

// AddChartRelationship adds a "relationships" reference to a chartVersion resource's chart
func AddChartRelationship(db datastore.Database, resource *swaggermodels.Resource, chartPackage *swaggermodels.ChartPackage) {
	resource.Relationships = &swaggermodels.ChartRelationship{
		Chart: &swaggermodels.ChartAsRelationship{
			Links: &swaggermodels.ResourceLink{
				Self: pointerto.String(MakeRepoChartRouteURL(APIVer1String, chartPackage.Repo, *chartPackage.Name)),
			},
			Data: &swaggermodels.Chart{
				Name:        chartPackage.Name,
				Description: chartPackage.Description,
				Repo:        getRepoObject(db, chartPackage),
				Home:        chartPackage.Home,
				Sources:     chartPackage.Sources,
				Maintainers: chartPackage.Maintainers,
			},
		},
	}
}

// AddLatestChartVersionRelationship adds a "relationships" reference to a chart resource's latest chartVersion
func AddLatestChartVersionRelationship(resource *swaggermodels.Resource, chartPackage *swaggermodels.ChartPackage) {
	resource.Relationships = &swaggermodels.LatestChartVersionRelationship{
		LatestChartVersion: &swaggermodels.ChartVersionAsRelationship{
			Links: &swaggermodels.ResourceLink{
				Self: pointerto.String(MakeRepoChartVersionRouteURL(APIVer1String, chartPackage.Repo, *chartPackage.Name, *chartPackage.Version)),
			},
			Data: &swaggermodels.ChartVersion{
				Created:    chartPackage.Created,
				Digest:     chartPackage.Digest,
				Urls:       chartPackage.Urls,
				Version:    chartPackage.Version,
				AppVersion: chartPackage.AppVersion,
				Icons:      makeAvailableIcons(chartPackage),
				Readme:     makeReadmeURL(chartPackage),
			},
		},
	}
}

// AddCanonicalLink adds a "self" link to a chart resource's canonical API endpoint
func AddCanonicalLink(resource *swaggermodels.Resource) {
	resource.Links = &swaggermodels.ResourceLink{
		Self: pointerto.String(MakeRepoChartRouteURL(APIVer1String, *resource.Attributes.(*swaggermodels.Chart).Repo.Name, *resource.Attributes.(*swaggermodels.Chart).Name)),
	}
}

// GetLatestChartVersion returns the most recent version from a slice of versioned charts
func GetLatestChartVersion(charts []*swaggermodels.ChartPackage, name string) (*swaggermodels.ChartPackage, error) {
	latest := "0.0.0"
	var ret *swaggermodels.ChartPackage
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
func GetChartVersion(charts []*swaggermodels.ChartPackage, name, version string) (*swaggermodels.ChartPackage, error) {
	var ret *swaggermodels.ChartPackage
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
func GetChartVersions(charts []*swaggermodels.ChartPackage, name string) ([]*swaggermodels.ChartPackage, error) {
	var ret []*swaggermodels.ChartPackage
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

func makeAvailableIcons(chart *swaggermodels.ChartPackage) []*swaggermodels.Icon {
	var res []*swaggermodels.Icon
	icons := charthelper.AvailableIcons(chart, "/assets")
	for _, icon := range icons {
		res = append(res, &swaggermodels.Icon{Name: &icon.Name, Path: &icon.Path})
	}
	return res
}

func makeReadmeURL(chart *swaggermodels.ChartPackage) *string {
	res := charthelper.ReadmeStaticUrl(chart, "/assets")
	return &res
}

func getRepoObject(db datastore.Database, chart *swaggermodels.ChartPackage) *swaggermodels.Repo {
	repos, err := models.ListRepos(db)
	if err != nil {
		log.Fatal("could not get Repo collection", err)
	}

	var repoPayload swaggermodels.Repo
	for _, repo := range repos {
		if repo.Name == chart.Repo {
			repoPayload = swaggermodels.Repo{
				Name:   &repo.Name,
				URL:    &repo.URL,
				Source: repo.Source,
			}
			return &repoPayload
		}
	}
	log.WithFields(log.Fields{"repo": chart.Repo, "chart": *chart.Name}).Error("could not find repo for chart")
	return &repoPayload
}
