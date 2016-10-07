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
	chartCreated     = "2016-05-26 11:23:44.086354411 +0000 UTC"
	chartChecksum    = "68eb4f96567c1d5fa9417b2bb9b9cbb2"
	chartDescription = "Chart for Apache HTTP Server"
	chartVersion     = "0.0.1"
	chartHome        = "https://github.com/kubernetes/charts/apache"
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
	assert.Equal(t, *charts[0].URL, chartURL, "chart URL field value")
	assert.Equal(t, *charts[0].Created, chartCreated, "chart created field value")
	assert.Equal(t, *charts[0].Checksum, chartChecksum, "chart checksum field value")
	assert.Equal(t, *charts[0].Description, chartDescription, "chart description field value")
	assert.Equal(t, *charts[0].Version, chartVersion, "chart version field value")
	assert.Equal(t, *charts[0].Home, chartHome, "chart home field value")
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
	assert.Equal(t, *chartResource.Links.(*models.ChartResourceLinks).Latest, fmt.Sprintf("/v1/charts/%s/%s/%s", repoName, chartName, chartVersion), "chart resource Links.Latest field value")
	assert.Equal(t, *chartResource.Attributes.(*models.ChartResourceAttributes).Repo, repoName, "chart resource Attributes.Repo field value")
	assert.Equal(t, *chartResource.Attributes.(*models.ChartResourceAttributes).Name, chartName, "chart resource Attributes.Name field value")
	assert.Equal(t, *chartResource.Attributes.(*models.ChartResourceAttributes).Description, chartDescription, "chart resource Attributes.Description field value")
	assert.Equal(t, *chartResource.Attributes.(*models.ChartResourceAttributes).Created, chartCreated, "chart resource Attributes.Created field value")
	assert.Equal(t, *chartResource.Attributes.(*models.ChartResourceAttributes).Home, chartHome, "chart resource Attributes.Home field value")
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
%s-%s:
  name: %s
  url: %s
  created: %s
  checksum: %s
  description: %s
  version: %s
  home: %s`, chartName, chartVersion, chartName, chartURL, chartCreated, chartChecksum, chartDescription, chartVersion, chartHome))
}
