package cache

import (
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/kubernetes-helm/monocular/src/api/config/repos"
)

func TestNewRefreshData(t *testing.T) {
	repos := repos.Repos{}
	chartsImplementation := NewCachedCharts(repos)
	// Setup background index refreshes
	freshness := time.Duration(3600) * time.Second
	job := NewRefreshChartsData(chartsImplementation, freshness, "test-run", false)
	assert.Equal(t, job.Frequency(), freshness, "Frequency")
	assert.Equal(t, job.FirstRun(), false, "First run")
	assert.Equal(t, job.Name(), "test-run", "Name")
	err := job.Do()
	assert.NoErr(t, err)
}

func TestNewRefreshDataError(t *testing.T) {
	repos := repos.Repos{
		repos.Repo{
			Name: "waps",
			URL:  "./localhost",
		},
	}
	chartsImplementation := NewCachedCharts(repos)
	freshness := time.Duration(3600) * time.Second
	job := NewRefreshChartsData(chartsImplementation, freshness, "test-run", true)
	assert.Equal(t, job.FirstRun(), true, "First run")
	err := job.Do()
	assert.ExistsErr(t, err, "Error executing refresh")
}
