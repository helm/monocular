package mocks

import (
	"testing"

	"github.com/arschles/assert"
)

const (
	repoName  = "stable"
	chartName = "apache"
)

func TestGetChartFromMockRepo(t *testing.T) {
	chart, err := GetChartFromMockRepo(repoName, chartName)
	assert.NoErr(t, err)
	assert.Equal(t, *chart.ID, repoName+"/"+chartName, "chart ID")
	chart, err = GetChartFromMockRepo("bogon", chartName)
	assert.ExistsErr(t, err, "sent bogus repo name to GetChartFromMockRepo")
	assert.Nil(t, chart.ID, "zero value ID")
}

func TestGetAllChartsFromMockRepos(t *testing.T) {
	charts, err := GetAllChartsFromMockRepos()
	assert.NoErr(t, err)
	assert.True(t, len(charts) > 0, "at least 1 chart returned")
}

func TestGetChartsFromMockRepo(t *testing.T) {
	charts, err := GetChartsFromMockRepo(repoName)
	assert.NoErr(t, err)
	assert.True(t, len(charts) > 0, "at least 1 chart returned")
	charts, err = GetChartsFromMockRepo("unparseable")
	assert.ExistsErr(t, err, "sent unparseable repo name to GetChartsFromMockRepo")
	assert.True(t, len(charts) == 0, "empty charts slice returned")
	charts, err = GetChartsFromMockRepo("bogon")
	assert.ExistsErr(t, err, "sent bogus repo name to GetChartsFromMockRepo")
	assert.True(t, len(charts) == 0, "empty charts slice returned")
}
