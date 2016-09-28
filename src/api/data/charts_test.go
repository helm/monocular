package data

import (
	"testing"

	"github.com/arschles/assert"
)

const (
	repoName  = "stable"
	chartName = "apache"
)

func TestGetChart(t *testing.T) {
	chart, err := GetChart(repoName, chartName)
	assert.NoErr(t, err)
	assert.Equal(t, *chart.ID, repoName+"/"+chartName, "chart ID")
	chart, err = GetChart("bogon", chartName)
	assert.ExistsErr(t, err, "sent bogus repo name to GetChart")
	assert.Nil(t, chart.ID, "zero value ID")
}

func TestGetAllCharts(t *testing.T) {
	_, err := GetAllCharts()
	assert.NoErr(t, err)
}

func TestGetChartsInRepo(t *testing.T) {
	charts, err := GetChartsInRepo(repoName)
	assert.NoErr(t, err)
	assert.True(t, len(charts) > 0, "returned charts")
	noCharts, err := GetChartsInRepo("bogon")
	assert.ExistsErr(t, err, "sent bogus repo name to GetChartsInRepo")
	assert.True(t, len(noCharts) == 0, "empty charts slice")
}
