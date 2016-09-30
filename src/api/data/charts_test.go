package data

import (
	"testing"

	"github.com/arschles/assert"
	"github.com/helm/monocular/src/api/pkg/testutil"
)

func TestGetChart(t *testing.T) {
	chart, err := GetChart(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	assert.Equal(t, *chart.ID, testutil.RepoName+"/"+testutil.ChartName, "chart ID")
	chart, err = GetChart(testutil.BogusRepo, testutil.ChartName)
	assert.ExistsErr(t, err, "sent bogus repo name to GetChart")
	assert.Nil(t, chart.ID, "zero value ID")
}

func TestGetAllCharts(t *testing.T) {
	_, err := GetAllCharts()
	assert.NoErr(t, err)
}

func TestGetChartsInRepo(t *testing.T) {
	charts, err := GetChartsInRepo(testutil.RepoName)
	assert.NoErr(t, err)
	assert.True(t, len(charts) > 0, "returned charts")
	noCharts, err := GetChartsInRepo("bogon")
	assert.ExistsErr(t, err, "sent bogus repo name to GetChartsInRepo")
	assert.True(t, len(noCharts) == 0, "empty charts slice")
}
