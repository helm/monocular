package cache

import (
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/kubernetes-helm/monocular/src/api/data/util"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
)

func TestNewRefreshData(t *testing.T) {
	setupTestRepoCache(nil)
	defer teardownTestRepoCache()

	chartsImplementation := NewCachedCharts()
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
	repos := []models.Repo{
		models.Repo{
			Name: util.StrToPtr("waps"),
			URL:  util.StrToPtr("./localhost"),
		},
	}
	setupTestRepoCache(&repos)
	defer teardownTestRepoCache()

	chartsImplementation := NewCachedCharts()
	freshness := time.Duration(3600) * time.Second
	job := NewRefreshChartsData(chartsImplementation, freshness, "test-run", true)
	assert.Equal(t, job.FirstRun(), true, "First run")
	err := job.Do()
	assert.ExistsErr(t, err, "Error executing refresh")
}
