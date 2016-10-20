package helpers

import (
	"fmt"
	"testing"

	"github.com/arschles/assert"
	"github.com/helm/monocular/src/api/swagger/models"
)

const (
	repoName         = "stable"
	chartName        = "apache"
	chartURL         = "https://storage.googleapis.com/kubernetes-charts/apache-0.0.1.tgz"
	chartSource      = "https://github.com/kubernetes/charts/apache"
	chartCreated     = "2016-10-06T16:23:20.499814565-06:00"
	chartDigest      = "99c76e403d752c84ead610644d4b1c2f2b453a74b921f422b9dcb8a7c8b559cd"
	chartDescription = "Chart for Apache HTTP Server"
	chartVersion     = "0.0.1"
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
	charts, err := ParseYAMLRepo(getTestRepoYAML())
	assert.NoErr(t, err)
	assert.Equal(t, len(charts), 1, "charts slice response from ParseYAMLRepo")
	assert.Equal(t, *charts[0].Name, chartName, "chart name field value")
	assert.Equal(t, *charts[0].Created, chartCreated, "chart created field value")
	assert.Equal(t, *charts[0].Digest, chartDigest, "chart checksum field value")
	assert.Equal(t, *charts[0].Description, chartDescription, "chart description field value")
	assert.Equal(t, *charts[0].Version, chartVersion, "chart version field value")
	assert.Equal(t, *charts[0].Home, chartHome, "chart home field value")
	//assert.Equal(t, *charts[0].Urls[0], chartURL, "chart URL field value")
	//assert.Equal(t, *charts[0].Sources[0], chartSource, "chart URL field value")
	_, err = ParseYAMLRepo([]byte(fmt.Sprintf(`this is not yaml`)))
	assert.ExistsErr(t, err, "sent something not yaml to ParseYAMLRepo")
	_, err = ParseYAMLRepo([]byte(fmt.Sprintf(`andy: kaufman`)))
	assert.ExistsErr(t, err, "sent bogus repo yaml ParseYAMLRepo")
}

func TestMakeChartResource(t *testing.T) {
	charts, err := ParseYAMLRepo(getTestRepoYAML())
	assert.NoErr(t, err)
	chartResource := MakeChartResource(charts[0], repoName)
	assert.Equal(t, *chartResource.Type, "chart", "chart resource type field value")
	assert.Equal(t, *chartResource.ID, repoName+"/"+chartName, "chart resource ID field value")
	assert.Equal(t, *chartResource.Attributes.(*models.ChartResourceAttributes).Repo, repoName, "chart resource Attributes.Repo field value")
	assert.Equal(t, *chartResource.Attributes.(*models.ChartResourceAttributes).Name, chartName, "chart resource Attributes.Name field value")
	assert.Equal(t, *chartResource.Attributes.(*models.ChartResourceAttributes).Description, chartDescription, "chart resource Attributes.Description field value")
	assert.Equal(t, *chartResource.Attributes.(*models.ChartResourceAttributes).Home, chartHome, "chart resource Attributes.Home field value")
}

func TestMakeChartsResource(t *testing.T) {
	charts, err := ParseYAMLRepo(getTestRepoYAML())
	assert.NoErr(t, err)
	chartsResource := MakeChartsResource(charts, repoName)
	assert.Equal(t, *chartsResource[0].Type, "chart", "chart resource type field value")
	assert.Equal(t, *chartsResource[0].ID, repoName+"/"+chartName, "chart resource ID field value")
	assert.Equal(t, *chartsResource[0].Attributes.(*models.ChartResourceAttributes).Repo, repoName, "chart resource Attributes.Repo field value")
	assert.Equal(t, *chartsResource[0].Attributes.(*models.ChartResourceAttributes).Name, chartName, "chart resource Attributes.Name field value")
	assert.Equal(t, *chartsResource[0].Attributes.(*models.ChartResourceAttributes).Description, chartDescription, "chart resource Attributes.Description field value")
	assert.Equal(t, *chartsResource[0].Attributes.(*models.ChartResourceAttributes).Home, chartHome, "chart resource Attributes.Home field value")
}

func TestMakeChartVersionResource(t *testing.T) {
	charts, err := ParseYAMLRepo(getTestRepoYAML())
	assert.NoErr(t, err)
	chartVersionResource := MakeChartVersionResource(charts[0], repoName)
	assert.Equal(t, *chartVersionResource.Type, "chartVersion", "chart resource type field value")
	assert.Equal(t, *chartVersionResource.ID, repoName+"/"+chartName+":"+chartVersion, "chart resource ID field value")
	assert.Equal(t, *chartVersionResource.Attributes.(*models.ChartVersionResourceAttributes).Repo, repoName, "chartVersion resource Attributes.Repo field value")
	assert.Equal(t, *chartVersionResource.Attributes.(*models.ChartVersionResourceAttributes).Name, chartName, "chartVersion resource Attributes.Name field value")
	assert.Equal(t, *chartVersionResource.Attributes.(*models.ChartVersionResourceAttributes).Version, chartVersion, "chartVersion resource Attributes.Version field value")
	assert.Equal(t, *chartVersionResource.Attributes.(*models.ChartVersionResourceAttributes).Description, chartDescription, "chartVersion resource Attributes.Description field value")
	assert.Equal(t, *chartVersionResource.Attributes.(*models.ChartVersionResourceAttributes).Home, chartHome, "chartVersion resource Attributes.Home field value")
}

func TestMakeChartVersionsResource(t *testing.T) {
	charts, err := ParseYAMLRepo(getTestRepoYAML())
	assert.NoErr(t, err)
	chartVersionsResource := MakeChartVersionsResource(charts, repoName)
	assert.Equal(t, *chartVersionsResource[0].Type, "chartVersion", "chart resource type field value")
	assert.Equal(t, *chartVersionsResource[0].ID, repoName+"/"+chartName+":"+chartVersion, "chart resource ID field value")
	assert.Equal(t, *chartVersionsResource[0].Attributes.(*models.ChartVersionResourceAttributes).Repo, repoName, "chartVersion resource Attributes.Repo field value")
	assert.Equal(t, *chartVersionsResource[0].Attributes.(*models.ChartVersionResourceAttributes).Name, chartName, "chartVersion resource Attributes.Name field value")
	assert.Equal(t, *chartVersionsResource[0].Attributes.(*models.ChartVersionResourceAttributes).Version, chartVersion, "chartVersion resource Attributes.Version field value")
	assert.Equal(t, *chartVersionsResource[0].Attributes.(*models.ChartVersionResourceAttributes).Description, chartDescription, "chartVersion resource Attributes.Description field value")
	assert.Equal(t, *chartVersionsResource[0].Attributes.(*models.ChartVersionResourceAttributes).Home, chartHome, "chartVersion resource Attributes.Home field value")
}

func TestAddCanonicalLink(t *testing.T) {
	charts, err := ParseYAMLRepo(getTestRepoYAML())
	assert.NoErr(t, err)
	chartResource := MakeChartResource(charts[0], repoName)
	AddCanonicalLink(chartResource)
	assert.Equal(t, *chartResource.Links.(*models.ChartResourceLinks).Canonical, fmt.Sprintf("/%s/charts/%s/%s", apiVer1, repoName, chartName), "chart resource Links.Latest field value")
}

func TestGetLatestChartVersion(t *testing.T) {
	charts, err := ParseYAMLRepo(getTestRepoYAML())
	assert.NoErr(t, err)
	moreCharts, err := ParseYAMLRepo(getTestRepoYAML())
	assert.NoErr(t, err)
	reallyLargeVersion := "999.0.0"
	*moreCharts[0].Version = reallyLargeVersion
	charts = append(charts, moreCharts[0])
	assert.Equal(t, len(charts), 2, "number of charts in charts array")
	latest, err := GetLatestChartVersion(charts, chartName)
	assert.NoErr(t, err)
	assert.Equal(t, *latest.Version, reallyLargeVersion, "latest chart version")
	chartsBadVersion, err := ParseYAMLRepo(getTestRepoYAML())
	assert.NoErr(t, err)
	*chartsBadVersion[0].Version = "this is not semver"
	latest, err = GetLatestChartVersion(chartsBadVersion, chartName)
	assert.ExistsErr(t, err, "sent chart with bogus version to GetLatestChartVersion")
}

func TestGetChartVersion(t *testing.T) {
	charts, err := ParseYAMLRepo(getTestRepoYAML())
	assert.NoErr(t, err)
	versionedCharts, err := GetChartVersion(charts, chartName, chartVersion)
	assert.NoErr(t, err)
	assert.Equal(t, *versionedCharts.Name, chartName, "chart name")
	assert.Equal(t, *versionedCharts.Version, chartVersion, "chart version")
	_, err = GetChartVersion(charts, chartName, "99.99.99")
	assert.ExistsErr(t, err, "requested non-existent version of chart")
}

func TestGetChartVersions(t *testing.T) {
	charts, err := ParseYAMLRepo(getTestRepoYAML())
	assert.NoErr(t, err)
	versionedCharts, err := GetChartVersions(charts, chartName)
	assert.NoErr(t, err)
	assert.Equal(t, *versionedCharts[0].Name, chartName, "chart name")
	_, err = GetChartVersions(charts, "noname")
	assert.ExistsErr(t, err, "requested versions of non-existent chart name")
}

func TestAddLatestRelationship(t *testing.T) {
	charts, err := ParseYAMLRepo(getTestRepoYAML())
	assert.NoErr(t, err)
	chartResource := MakeChartResource(charts[0], repoName)
	AddLatestRelationship(chartResource, charts[0])
	// TODO validate data
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

func TestInt64ToPtr(t *testing.T) {
	var number int64
	number = 13
	ptr := Int64ToPtr(number)
	assert.Equal(t, number, *ptr, "int64 to ptr conversion")
}

func TestStrToPtr(t *testing.T) {
	var str string
	str = "string"
	ptr := StrToPtr(str)
	assert.Equal(t, str, *ptr, "string to ptr conversion")
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
generated: 2016-10-06T16:23:20.499029981-06:00`, apiVer1, chartCreated, chartDescription, chartDigest, chartHome, chartName, chartSource, chartURL, chartVersion))
}
