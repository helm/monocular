package mocks

import (
	"testing"

	"github.com/arschles/assert"
	"github.com/helm/monocular/src/api/pkg/testutil"
)

var chartsImplementation = NewMockCharts()

func TestMockChartsChartFromRepo(t *testing.T) {
	chart, err := chartsImplementation.ChartFromRepo(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	assert.Equal(t, *chart.ID, testutil.RepoName+"/"+testutil.ChartName, "chart ID")
	chart, err = chartsImplementation.ChartFromRepo(testutil.BogusRepo, testutil.ChartName)
	assert.ExistsErr(t, err, "sent bogus repo name to Charts.ChartFromRepo()")
	assert.Nil(t, chart.ID, "zero value ID")
}

func TestMockChartsAll(t *testing.T) {
	_, err := chartsImplementation.All()
	assert.NoErr(t, err)
}

func TestMockChartsAllFromRepo(t *testing.T) {
	charts, err := chartsImplementation.AllFromRepo(testutil.RepoName)
	assert.NoErr(t, err)
	assert.True(t, len(charts) > 0, "returned charts")
	noCharts, err := chartsImplementation.AllFromRepo(testutil.BogusRepo)
	assert.ExistsErr(t, err, "sent bogus repo name to GetChartsInRepo")
	assert.True(t, len(noCharts) == 0, "empty charts slice")
}

func TestGetChartFromMockRepo(t *testing.T) {
	chart, err := GetChartFromMockRepo(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	assert.Equal(t, *chart.ID, testutil.RepoName+"/"+testutil.ChartName, "chart ID")
	chart, err = GetChartFromMockRepo(testutil.BogusRepo, testutil.ChartName)
	assert.ExistsErr(t, err, "sent bogus repo name to GetChartFromMockRepo")
	assert.Nil(t, chart.ID, "zero value ID")
	chart, err = GetChartFromMockRepo("unparseable", testutil.ChartName)
	assert.ExistsErr(t, err, "sent unparseable repo name to GetChartsFromMockRepo")
	assert.Nil(t, chart.ID, "zero value ID")
}

func TestGetAllChartsFromMockRepos(t *testing.T) {
	charts, err := GetAllChartsFromMockRepos()
	assert.NoErr(t, err)
	assert.True(t, len(charts) > 0, "at least 1 chart returned")
}

func TestGetChartsFromMockRepo(t *testing.T) {
	charts, err := GetChartsFromMockRepo(testutil.RepoName)
	assert.NoErr(t, err)
	assert.True(t, len(charts) > 0, "at least 1 chart returned")
	charts, err = GetChartsFromMockRepo("unparseable")
	assert.ExistsErr(t, err, "sent unparseable repo name to GetChartsFromMockRepo")
	assert.True(t, len(charts) == 0, "empty charts slice returned")
	charts, err = GetChartsFromMockRepo(testutil.BogusRepo)
	assert.ExistsErr(t, err, "sent bogus repo name to GetChartsFromMockRepo")
	assert.True(t, len(charts) == 0, "empty charts slice returned")
}
