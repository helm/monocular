package cache

import (
	"testing"

	"github.com/arschles/assert"
	"github.com/helm/monocular/src/api/testutil"
)

var repos = []map[string]string{
	map[string]string{
		"stable": "https://github.com/kubernetes/charts",
	},
	map[string]string{
		"incubator": "https://github.com/kubernetes/charts/tree/master/incubator",
	},
}
var chartsImplementation = NewCachedCharts(repos)

func TestCachedChartsChartFromRepo(t *testing.T) {
	err := chartsImplementation.Refresh()
	assert.NoErr(t, err)
	// TODO: validate chart data
	_, err = chartsImplementation.ChartFromRepo(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	_, err = chartsImplementation.ChartFromRepo(testutil.BogusRepo, testutil.ChartName)
	assert.ExistsErr(t, err, "sent bogus repo name to Charts.ChartFromRepo()")
	_, err = chartsImplementation.ChartFromRepo(testutil.RepoName, testutil.BogusRepo)
	assert.ExistsErr(t, err, "sent bogus chart name to Charts.ChartFromRepo()")
}

func TestCachedChartsAll(t *testing.T) {
	err := chartsImplementation.Refresh()
	assert.NoErr(t, err)
	_, err = chartsImplementation.All()
	assert.NoErr(t, err)
}

func TestCachedChartsAllFromRepo(t *testing.T) {
	err := chartsImplementation.Refresh()
	assert.NoErr(t, err)
	charts, err := chartsImplementation.AllFromRepo(testutil.RepoName)
	assert.NoErr(t, err)
	assert.True(t, len(charts) > 0, "returned charts")
	noCharts, err := chartsImplementation.AllFromRepo(testutil.BogusRepo)
	assert.ExistsErr(t, err, "sent bogus repo name to GetChartsInRepo")
	assert.True(t, len(noCharts) == 0, "empty charts slice")
}

func TestCachedChartsRefresh(t *testing.T) {
	err := chartsImplementation.Refresh()
	assert.NoErr(t, err)
}
