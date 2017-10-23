package mocks

import (
	"errors"
	"testing"

	"github.com/arschles/assert"
	"github.com/kubernetes-helm/monocular/src/api/data/helpers"
	"github.com/kubernetes-helm/monocular/src/api/models"
	swaggermodels "github.com/kubernetes-helm/monocular/src/api/swagger/models"
	"github.com/kubernetes-helm/monocular/src/api/swagger/restapi/operations/charts"
	"github.com/kubernetes-helm/monocular/src/api/testutil"
)

var chartsImplementation = NewMockCharts(MockedMethods{})

func TestMockChartsChartFromRepo(t *testing.T) {
	// TODO: validate chart data
	_, err := chartsImplementation.ChartFromRepo(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	_, err = chartsImplementation.ChartFromRepo(testutil.BogusRepo, testutil.ChartName)
	assert.ExistsErr(t, err, "sent bogus repo name to Charts.ChartFromRepo()")
	_, err = chartsImplementation.ChartFromRepo(testutil.RepoName, testutil.BogusRepo)
	assert.ExistsErr(t, err, "sent bogus chart name to Charts.ChartFromRepo()")
}

func TestMockChartsChartVersionFromRepo(t *testing.T) {
	chart, err := chartsImplementation.ChartVersionFromRepo(testutil.RepoName, testutil.ChartName, testutil.ChartVersionString)
	assert.NoErr(t, err)
	assert.Equal(t, *chart.Name, testutil.ChartName, "chart name")
	assert.Equal(t, *chart.Version, testutil.ChartVersionString, "chart version")
	_, err = chartsImplementation.ChartVersionFromRepo(testutil.RepoName, testutil.ChartName, "99.99.99")
	assert.ExistsErr(t, err, "sent bogus chart version to ChartVersionFromRepo")
	_, err = chartsImplementation.ChartVersionFromRepo(testutil.BogusRepo, testutil.ChartName, testutil.ChartVersionString)
	assert.ExistsErr(t, err, "sent bogus repo name to Charts.ChartFromRepo()")
	_, err = chartsImplementation.ChartVersionFromRepo(testutil.RepoName, testutil.BogusRepo, testutil.ChartVersionString)
	assert.ExistsErr(t, err, "sent bogus chart name to Charts.ChartFromRepo()")
	_, err = chartsImplementation.ChartVersionFromRepo(testutil.UnparseableRepo, testutil.ChartName, testutil.ChartVersionString)
	assert.ExistsErr(t, err, "sent unparseable repo name to ChartVersionFromRepo")
}

func TestMockChartsChartVersionsFromRepo(t *testing.T) {
	charts, err := chartsImplementation.ChartVersionsFromRepo(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	assert.True(t, len(charts) > 0, "returned charts")
	noCharts, err := chartsImplementation.ChartVersionsFromRepo(testutil.BogusRepo, testutil.ChartName)
	assert.ExistsErr(t, err, "sent bogus repo name to GetChartsInRepo")
	assert.True(t, len(noCharts) == 0, "empty charts slice")
	noCharts, err = chartsImplementation.ChartVersionsFromRepo(testutil.UnparseableRepo, testutil.ChartName)
	assert.ExistsErr(t, err, "sent unparseable repo name to GetChartsInRepo")
	assert.True(t, len(noCharts) == 0, "empty charts slice")
	noCharts, err = chartsImplementation.ChartVersionsFromRepo(testutil.RepoName, testutil.BogusRepo)
	assert.ExistsErr(t, err, "sent bogus chart name to GetChartsInRepo")
	assert.True(t, len(noCharts) == 0, "empty charts slice")
}

func TestMockChartsAll(t *testing.T) {
	_, err := chartsImplementation.All()
	assert.NoErr(t, err)
}

func TestMockChartsAllWithMockedMethod(t *testing.T) {
	chImplementation := NewMockCharts(MockedMethods{
		All: func() ([]*swaggermodels.ChartPackage, error) {
			var ret []*swaggermodels.ChartPackage
			return ret, errors.New("error getting all charts")
		},
	})
	_, err := chImplementation.All()
	assert.ExistsErr(t, err, "mocked error")
}

func TestMockChartsSearch(t *testing.T) {
	params := charts.SearchChartsParams{
		Name: "drupal",
	}
	charts, err := chartsImplementation.Search(params)
	assert.NoErr(t, err)
	// flatten chart+version results into a chart resource array
	db, _ := models.NewMockSession(models.MockDBConfig{}).DB()
	resources := helpers.MakeChartResources(db, charts)
	assert.Equal(t, len(resources), 1, "number of unique chart results")
}

func TestMockChartsSearchWithMockedMethod(t *testing.T) {
	chImplementation := NewMockCharts(MockedMethods{
		Search: func(params charts.SearchChartsParams) ([]*swaggermodels.ChartPackage, error) {
			var ret []*swaggermodels.ChartPackage
			return ret, errors.New("error searching charts")
		},
	})
	params := charts.SearchChartsParams{
		Name: "drupal",
	}
	_, err := chImplementation.Search(params)
	assert.ExistsErr(t, err, "mocked error")
}

func TestMockChartsAllFromRepo(t *testing.T) {
	charts, err := chartsImplementation.AllFromRepo(testutil.RepoName)
	assert.NoErr(t, err)
	assert.True(t, len(charts) > 0, "returned charts")
	noCharts, err := chartsImplementation.AllFromRepo(testutil.BogusRepo)
	assert.ExistsErr(t, err, "sent bogus repo name to GetChartsInRepo")
	assert.True(t, len(noCharts) == 0, "empty charts slice")
}
