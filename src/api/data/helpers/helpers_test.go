package helpers

import (
	"fmt"
	"testing"

	"github.com/arschles/assert"
	"github.com/kubernetes-helm/monocular/src/api/data/cache/charthelper"
	"github.com/kubernetes-helm/monocular/src/api/datastore"
	"github.com/kubernetes-helm/monocular/src/api/models"
	swaggermodels "github.com/kubernetes-helm/monocular/src/api/swagger/models"
)

const (
	repoName         = "testRepo"
	chartName        = "apache"
	chartURL         = "https://storage.googleapis.com/kubernetes-charts/apache-0.0.1.tgz"
	chartSource      = "https://github.com/kubernetes/charts/apache"
	chartCreated     = "2016-10-06T16:23:20.499814565-06:00"
	chartDigest      = "99c76e403d752c84ead610644d4b1c2f2b453a74b921f422b9dcb8a7c8b559cd"
	chartDescription = "Chart for Apache HTTP Server"
	chartVersion     = "0.0.1"
	chartAppVersion  = "1.0.0"
	chartHome        = "https://k8s.io/helm"
)

func TestIsYAML(t *testing.T) {
	yaml := []byte(fmt.Sprintf(`
root-property:
  sub-property: value`))
	assert.Equal(t, IsYAML(yaml), true, "YAML to IsYAML helper function")
	byteStr := []byte(fmt.Sprintf(`this is a string`))
	assert.Equal(t, IsYAML(byteStr), false, "string to IsYAML helper function")
}

func TestParseYAMLRepo(t *testing.T) {
	charts, err := ParseYAMLRepo(getTestRepoYAML(), repoName)
	assert.NoErr(t, err)
	assert.Equal(t, len(charts), 1, "charts slice response from ParseYAMLRepo")
	assert.Equal(t, *charts[0].Name, chartName, "chart name field value")
	assert.Equal(t, *charts[0].Created, chartCreated, "chart created field value")
	assert.Equal(t, *charts[0].Digest, chartDigest, "chart checksum field value")
	assert.Equal(t, *charts[0].Description, chartDescription, "chart description field value")
	assert.Equal(t, *charts[0].Version, chartVersion, "chart version field value")
	assert.Equal(t, *charts[0].AppVersion, chartAppVersion, "chart app version field value")
	assert.Equal(t, *charts[0].Home, chartHome, "chart home field value")
	//assert.Equal(t, *charts[0].Urls[0], chartURL, "chart URL field value")
	//assert.Equal(t, *charts[0].Sources[0], chartSource, "chart URL field value")
	_, err = ParseYAMLRepo([]byte(fmt.Sprintf(`this is not yaml`)), repoName)
	assert.ExistsErr(t, err, "sent something not yaml to ParseYAMLRepo")
	_, err = ParseYAMLRepo([]byte(fmt.Sprintf(`andy: kaufman`)), repoName)
	assert.ExistsErr(t, err, "sent bogus repo yaml ParseYAMLRepo")
}

func TestParseYAMLRepoWithDeprecatedChart(t *testing.T) {
	charts, err := ParseYAMLRepo(getTestRepoYAMLWithDeprecatedChart(), repoName)
	assert.NoErr(t, err)
	assert.Equal(t, len(charts), 1, "charts slice response from ParseYAMLRepo")
	assert.Equal(t, *charts[0].Name, chartName, "chart name field value")
	assert.Equal(t, *charts[0].Created, chartCreated, "chart created field value")
	assert.Equal(t, *charts[0].Digest, chartDigest, "chart checksum field value")
	assert.Equal(t, *charts[0].Description, chartDescription, "chart description field value")
	assert.Equal(t, *charts[0].Version, chartVersion, "chart version field value")
	assert.Equal(t, *charts[0].AppVersion, chartAppVersion, "chart app version field value")
	assert.Equal(t, *charts[0].Home, chartHome, "chart home field value")
}

func TestMakeChartResource(t *testing.T) {
	db := getTestDB()
	repo := models.OfficialRepos[0]
	repo.Name = repoName
	charts, err := ParseYAMLRepo(getTestRepoYAML(), repoName)
	assert.NoErr(t, err)
	assert.NoErr(t, err)
	chartResource := MakeChartResource(db, charts[0])
	assert.Equal(t, *chartResource.Type, "chart", "chart resource type field value")
	assert.Equal(t, *chartResource.ID, MakeChartID(repoName, chartName), "chart resource ID field value")
	assert.Equal(t, *chartResource.Attributes.(*swaggermodels.Chart).Repo.Name, repo.Name, "chart resource Attributes.Repo Name field value")
	assert.Equal(t, *chartResource.Attributes.(*swaggermodels.Chart).Repo.URL, repo.URL, "chart resource Attributes.Repo URL field value")
	assert.Equal(t, chartResource.Attributes.(*swaggermodels.Chart).Repo.Source, repo.Source, "chart resource Attributes.URL field value")
	assert.Equal(t, *chartResource.Attributes.(*swaggermodels.Chart).Name, chartName, "chart resource Attributes.Name field value")
	assert.Equal(t, *chartResource.Attributes.(*swaggermodels.Chart).Description, chartDescription, "chart resource Attributes.Description field value")
	assert.Equal(t, *chartResource.Attributes.(*swaggermodels.Chart).Home, chartHome, "chart resource Attributes.Home field value")
}

func TestMakeChartResources(t *testing.T) {
	db := getTestDB()
	charts, err := ParseYAMLRepo(getTestRepoYAML(), repoName)
	assert.NoErr(t, err)
	chartsResource := MakeChartResources(db, charts)
	repo := getRepoObject(db, charts[0])
	assert.Equal(t, *chartsResource[0].Type, "chart", "chart resource type field value")
	assert.Equal(t, *chartsResource[0].ID, MakeChartID(repoName, chartName), "chart resource ID field value")
	assert.Equal(t, *chartsResource[0].Attributes.(*swaggermodels.Chart).Repo.Name, *repo.Name, "chart resource Attributes.Repo Namefield value")
	assert.Equal(t, *chartsResource[0].Attributes.(*swaggermodels.Chart).Repo.URL, *repo.URL, "chart resource Attributes.Repo URL value")
	assert.Equal(t, chartsResource[0].Attributes.(*swaggermodels.Chart).Repo.Source, repo.Source, "chart resource Attributes.Source value")
	assert.Equal(t, *chartsResource[0].Attributes.(*swaggermodels.Chart).Name, chartName, "chart resource Attributes.Name field value")
	assert.Equal(t, *chartsResource[0].Attributes.(*swaggermodels.Chart).Description, chartDescription, "chart resource Attributes.Description field value")
	assert.Equal(t, *chartsResource[0].Attributes.(*swaggermodels.Chart).Home, chartHome, "chart resource Attributes.Home field value")
}

func TestMakeRepoResource(t *testing.T) {
	repo := models.OfficialRepos[0]
	repoResource := MakeRepoResource(repo)
	assert.Equal(t, *repoResource.Type, "repository", "repo resource type field value")
	assert.Equal(t, *repoResource.ID, repo.Name, "repo resource ID field value")
	assert.Equal(t, repoResource.Attributes.(*models.Repo).Name, repo.Name, "repo name")
	assert.Equal(t, repoResource.Attributes.(*models.Repo).URL, repo.URL, "repo URL")
	assert.Equal(t, repoResource.Attributes.(*models.Repo).Source, repo.Source, "chart resource Attributes.URL field value")
}

func TestMakeRepoResources(t *testing.T) {
	repos := models.OfficialRepos
	repoResource := MakeRepoResources(repos)[0]
	assert.Equal(t, *repoResource.Type, "repository", "repo resource type field value")
	assert.Equal(t, *repoResource.ID, repos[0].Name, "repo resource ID field value")
	assert.Equal(t, repoResource.Attributes.(*models.Repo).Name, repos[0].Name, "repo name")
	assert.Equal(t, repoResource.Attributes.(*models.Repo).URL, repos[0].URL, "repo URL")
	assert.Equal(t, repoResource.Attributes.(*models.Repo).Source, repos[0].Source, "chart resource Attributes.URL field value")
}

func TestMakeChartVersionResource(t *testing.T) {
	db := getTestDB()
	charts, err := ParseYAMLRepo(getTestRepoYAML(), repoName)
	assert.NoErr(t, err)
	chartVersionResource := MakeChartVersionResource(db, charts[0])
	assert.Equal(t, *chartVersionResource.Type, "chartVersion", "chart resource type field value")
	assert.Equal(t, *chartVersionResource.ID, MakeChartVersionID(repoName, chartName, chartVersion), "chart resource ID field value")
	assert.Equal(t, *chartVersionResource.Attributes.(*swaggermodels.ChartVersion).Created, chartCreated, "chartVersion resource Attributes.Created field value")
	assert.Equal(t, *chartVersionResource.Attributes.(*swaggermodels.ChartVersion).Digest, chartDigest, "chartVersion resource Attributes.Digest field value")
	assert.Equal(t, chartVersionResource.Attributes.(*swaggermodels.ChartVersion).Urls[0], chartURL, "chartVersion resource Attributes.Urls field value")
	assert.Equal(t, *chartVersionResource.Attributes.(*swaggermodels.ChartVersion).Version, chartVersion, "chartVersion resource Attributes.Version field value")
}

func TestMakeChartVersionResources(t *testing.T) {
	db := getTestDB()
	charts, err := ParseYAMLRepo(getTestRepoYAML(), repoName)
	assert.NoErr(t, err)
	chartVersionsResource := MakeChartVersionResources(db, charts)
	assert.Equal(t, *chartVersionsResource[0].Type, "chartVersion", "chart resource type field value")
	assert.Equal(t, *chartVersionsResource[0].ID, MakeChartVersionID(repoName, chartName, chartVersion), "chart resource ID field value")
	assert.Equal(t, *chartVersionsResource[0].Attributes.(*swaggermodels.ChartVersion).Created, chartCreated, "chartVersion resource Attributes.Created field value")
	assert.Equal(t, *chartVersionsResource[0].Attributes.(*swaggermodels.ChartVersion).Digest, chartDigest, "chartVersion resource Attributes.Digest field value")
	assert.Equal(t, chartVersionsResource[0].Attributes.(*swaggermodels.ChartVersion).Urls[0], chartURL, "chartVersion resource Attributes.Urls field value")
	assert.Equal(t, *chartVersionsResource[0].Attributes.(*swaggermodels.ChartVersion).Version, chartVersion, "chartVersion resource Attributes.Version field value")
}

func TestAddChartRelationship(t *testing.T) {
	db := getTestDB()
	charts, err := ParseYAMLRepo(getTestRepoYAML(), repoName)
	assert.NoErr(t, err)
	chart := charts[0]
	chartVersionResource := MakeChartVersionResource(db, chart)
	AddChartRelationship(db, chartVersionResource, chart)
	assert.Equal(t, *chartVersionResource.Relationships.(*swaggermodels.ChartRelationship).Chart.Data.Name, *chart.Name, "relationships.chart.data.name field value")
	assert.Equal(t, *chartVersionResource.Relationships.(*swaggermodels.ChartRelationship).Chart.Data.Description, *chart.Description, "relationships.chart.data.description field value")
	assert.Equal(t, *chartVersionResource.Relationships.(*swaggermodels.ChartRelationship).Chart.Data.Home, *chart.Home, "relationships.chart.data.home field value")
	assert.Equal(t, chartVersionResource.Relationships.(*swaggermodels.ChartRelationship).Chart.Data.Maintainers, chart.Maintainers, "relationships.chart.data.maintainers array value")
	assert.Equal(t, chartVersionResource.Relationships.(*swaggermodels.ChartRelationship).Chart.Data.Sources, chart.Sources, "relationships.chart.data.sources array value")
	assert.Equal(t, chartVersionResource.Relationships.(*swaggermodels.ChartRelationship).Chart.Data.Repo, getRepoObject(db, chart), "relationships.chart.data.repo field value")
	assert.Equal(t, *chartVersionResource.Relationships.(*swaggermodels.ChartRelationship).Chart.Links.Self, MakeRepoChartRouteURL(APIVer1String, chart.Repo, *chart.Name), "relationships.chart.links.self field value")
}

func TestAddLatestChartVersionRelationship(t *testing.T) {
	db := getTestDB()
	charts, err := ParseYAMLRepo(getTestRepoYAML(), repoName)
	assert.NoErr(t, err)
	chart := charts[0]
	chartResource := MakeChartResource(db, chart)
	AddLatestChartVersionRelationship(chartResource, chart)
	assert.Equal(t, *chartResource.Relationships.(*swaggermodels.LatestChartVersionRelationship).LatestChartVersion.Data.Created, *chart.Created, "relationships.latestChartVersion.data.created field value")
	assert.Equal(t, *chartResource.Relationships.(*swaggermodels.LatestChartVersionRelationship).LatestChartVersion.Data.Digest, *chart.Digest, "relationships.latestChartVersion.data.digest field value")
	assert.Equal(t, chartResource.Relationships.(*swaggermodels.LatestChartVersionRelationship).LatestChartVersion.Data.Urls, chart.Urls, "relationships.latestChartVersion.data.Urls field value")
	assert.Equal(t, *chartResource.Relationships.(*swaggermodels.LatestChartVersionRelationship).LatestChartVersion.Data.Version, *chart.Version, "relationships.latestChartVersion.data.digest field value")
	assert.Equal(t, *chartResource.Relationships.(*swaggermodels.LatestChartVersionRelationship).LatestChartVersion.Links.Self, MakeRepoChartVersionRouteURL(APIVer1String, chart.Repo, *chart.Name, *chart.Version), "relationships.chartVersion.links.self field value")
}

func TestAddCanonicalLink(t *testing.T) {
	db := getTestDB()
	charts, err := ParseYAMLRepo(getTestRepoYAML(), repoName)
	assert.NoErr(t, err)
	chartResource := MakeChartResource(db, charts[0])
	AddCanonicalLink(chartResource)
	assert.Equal(t, *chartResource.Links.(*swaggermodels.ResourceLink).Self, MakeRepoChartRouteURL(APIVer1String, repoName, chartName), "chart resource Links.Self field value")
}

func TestGetLatestChartVersion(t *testing.T) {
	charts, err := ParseYAMLRepo(getTestRepoYAML(), repoName)
	assert.NoErr(t, err)
	moreCharts, err := ParseYAMLRepo(getTestRepoYAML(), repoName)
	assert.NoErr(t, err)
	reallyLargeVersion := "999.0.0"
	*moreCharts[0].Version = reallyLargeVersion
	charts = append(charts, moreCharts[0])
	assert.Equal(t, len(charts), 2, "number of charts in charts array")
	latest, err := GetLatestChartVersion(charts, chartName)
	assert.NoErr(t, err)
	assert.Equal(t, *latest.Version, reallyLargeVersion, "latest chart version")
	chartsBadVersion, err := ParseYAMLRepo(getTestRepoYAML(), repoName)
	assert.NoErr(t, err)
	*chartsBadVersion[0].Version = "this is not semver"
	latest, err = GetLatestChartVersion(chartsBadVersion, chartName)
	assert.ExistsErr(t, err, "sent chart with bogus version to GetLatestChartVersion")
	latest, err = GetLatestChartVersion(charts, "no name")
	assert.ExistsErr(t, err, "Chart not found")
}

func TestGetChartVersion(t *testing.T) {
	charts, err := ParseYAMLRepo(getTestRepoYAML(), repoName)
	assert.NoErr(t, err)
	versionedCharts, err := GetChartVersion(charts, chartName, chartVersion)
	assert.NoErr(t, err)
	assert.Equal(t, *versionedCharts.Name, chartName, "chart name")
	assert.Equal(t, *versionedCharts.Version, chartVersion, "chart version")
	_, err = GetChartVersion(charts, chartName, "99.99.99")
	assert.ExistsErr(t, err, "requested non-existent version of chart")
}

func TestGetChartVersions(t *testing.T) {
	charts, err := ParseYAMLRepo(getTestRepoYAML(), repoName)
	assert.NoErr(t, err)
	versionedCharts, err := GetChartVersions(charts, chartName)
	assert.NoErr(t, err)
	assert.Equal(t, *versionedCharts[0].Name, chartName, "chart name")
	_, err = GetChartVersions(charts, "noname")
	assert.ExistsErr(t, err, "requested versions of non-existent chart name")
}

func TestNewestSemVer(t *testing.T) {
	// Verify that NewestSemVer returns correct semver string for larger major, minor, and patch substrings
	const v1Lower = "2.0.0"
	v2s := [3]string{"3.0.0", "2.1.0", "2.0.1"}
	for _, v2 := range v2s {
		newest, err := newestSemVer(v1Lower, v2)
		assert.NoErr(t, err)
		assert.Equal(t, v2, newest, "semver comparison")
	}
	// Verify that NewestSemVer returns correct semver string for smaller major, minor, and patch substrings
	const v1Higher = "2.4.5"
	v2s = [3]string{"1.99.23", "2.3.99", "2.4.4"}
	for _, v2 := range v2s {
		newest, err := newestSemVer(v1Higher, v2)
		assert.NoErr(t, err)
		assert.Equal(t, v1Higher, newest, "semver comparison")
	}
	// Verify that NewestSemVer returns correct semver string for comparing equal strings
	const v1Equal = "1.0.0"
	v2 := v1Equal
	newest, err := newestSemVer(v1Equal, v2)
	assert.NoErr(t, err)
	if newest != v1Equal && newest != v2 {
		fmt.Printf("expected %s to be equal to %s and %s\n", newest, v1Equal, v2)
		t.Fatal("semver comparison failure")
	}
	// Verify error conditions
	newest, err = newestSemVer("this is bogus", "and so is this")
	assert.ExistsErr(t, err, "sent bogus versions to newestSemVer")
	assert.Equal(t, newest, "", "newestSemVer response should be an empty string in an error case")
	newest, err = newestSemVer("1.0.0", "this is bogus")
	assert.ExistsErr(t, err, "sent bogus version as 1st arg to newestSemVer")
	assert.Equal(t, newest, "", "newestSemVer response should be an empty string in an error case")
	newest, err = newestSemVer("this is bogus", "1.0.0")
	assert.ExistsErr(t, err, "sent bogus version as 2nd arg to newestSemVer")
	assert.Equal(t, newest, "", "newestSemVer response should be an empty string in an error case")
}

func getTestRepoYAML() []byte {
	return []byte(fmt.Sprintf(`
apiVersion: %s
entries:
  apache:
    - created: %s
      description: %s
      digest: %s
      home: %s
      name: %s
      sources:
        - %s
      urls:
        - %s
      version: %s
      appVersion: %s
generated: 2016-10-06T16:23:20.499029981-06:00`, APIVer1String, chartCreated, chartDescription, chartDigest, chartHome, chartName, chartSource, chartURL, chartVersion, chartAppVersion))
}

func getTestRepoYAMLWithDeprecatedChart() []byte {
	return []byte(fmt.Sprintf(`
apiVersion: %s
entries:
  apache:
    - created: %s
      description: %s
      digest: %s
      home: %s
      name: %s
      sources:
        - %s
      urls:
        - %s
      version: %s
      appVersion: %s
  deprecated:
    - created: %s
      deprecated: true
      description: %s
      digest: %s
      home: %s
      name: deprecated
      sources:
        - %s
      urls:
        - %s
      version: 1.0.0
      appVersion: %s
    - created: %s
      description: %s
      digest: %s
      home: %s
      name: deprecated
      sources:
        - %s
      urls:
        - %s
      version: %s
      appVersion: %s
generated: 2016-10-06T16:23:20.499029981-06:00`,
		APIVer1String, chartCreated, chartDescription, chartDigest, chartHome, chartName, chartSource, chartURL, chartVersion, chartAppVersion,
		chartCreated, chartDescription, chartDigest, chartHome, chartSource, chartURL, chartAppVersion,
		chartCreated, chartDescription, chartDigest, chartHome, chartSource, chartURL, chartVersion, chartAppVersion,
	))
}

func TestMakeAvailableIcons(t *testing.T) {
	charts, err := ParseYAMLRepo(getTestRepoYAML(), repoName)
	assert.NoErr(t, err)
	chart := charts[0]
	AvailableIconsOrig := charthelper.AvailableIcons
	defer func() { charthelper.AvailableIcons = AvailableIconsOrig }()
	charthelper.AvailableIcons = func(chart *swaggermodels.ChartPackage, prefix string) []*charthelper.IconOutput {
		return []*charthelper.IconOutput{
			{"format1", "/myPath1"},
			{"format2", "mypath2"},
		}
	}
	iconOutputs := makeAvailableIcons(chart)

	for i, icon := range charthelper.AvailableIcons(chart, "prefix") {
		assert.Equal(t, *iconOutputs[i].Name, icon.Name, "Same name")
		assert.Equal(t, *iconOutputs[i].Path, icon.Path, "Same path")
	}
}

func TestGetRepoObject(t *testing.T) {
	db := getTestDB()
	charts, err := ParseYAMLRepo(getTestRepoYAML(), repoName)
	assert.NoErr(t, err)
	chart := charts[0]
	repo := getRepoObject(db, chart)
	assert.Equal(t, repo.Name, &chart.Repo, "Same repo Name")

	// Returns empty Repo if does not exist
	chart.Repo = "does-not-exist"
	repo = getRepoObject(db, chart)
	if repo.Name != nil || repo.URL != nil {
		t.Errorf("Repo Name and URL should be nil")
	}
}

func getTestDB() datastore.Database {
	repo := models.OfficialRepos[0]
	repo.Name = repoName
	db, _ := models.NewMockSession(models.MockDBConfig{}).DB()
	return db
}
