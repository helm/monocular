package cache

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/kubernetes-helm/monocular/src/api/config"
	"github.com/kubernetes-helm/monocular/src/api/data/pointerto"
	"github.com/kubernetes-helm/monocular/src/api/storage"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
)

func TestMain(m *testing.M) {
	flag.Parse()
	storage.Init(config.StorageConfig{"redis", ""})
	os.Exit(m.Run())
}

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
		{
			Name: pointerto.String("waps"),
			URL:  pointerto.String("./localhost"),
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
