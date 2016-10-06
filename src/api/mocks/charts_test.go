package mocks

import (
	"testing"

	"github.com/arschles/assert"
	"github.com/helm/monocular/src/api/pkg/testutil"
)

var chartsImplementation = NewMockCharts()

func TestMockChartsChartFromRepo(t *testing.T) {
	// TODO: validate chart data
	_, err := chartsImplementation.ChartFromRepo(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	_, err = chartsImplementation.ChartFromRepo(testutil.BogusRepo, testutil.ChartName)
	assert.ExistsErr(t, err, "sent bogus repo name to Charts.ChartFromRepo()")
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
	// TODO: validate chart data
	_, err := GetChartFromMockRepo(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	_, err = GetChartFromMockRepo(testutil.BogusRepo, testutil.ChartName)
	assert.ExistsErr(t, err, "sent bogus repo name to GetChartFromMockRepo")
	_, err = GetChartFromMockRepo("unparseable", testutil.ChartName)
	assert.ExistsErr(t, err, "sent unparseable repo name to GetChartsFromMockRepo")
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
