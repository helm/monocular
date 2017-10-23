package cache

import (
	"testing"
	"time"

	"github.com/arschles/assert"
	"github.com/kubernetes-helm/monocular/src/api/models"
)

func TestNewRefreshData(t *testing.T) {
	chartsImplementation := NewCachedCharts(dbSession)
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
	models.MockRepos = []*models.Repo{
		{
			Name: "waps",
			URL:  "./localhost",
		},
	}
	defer func() { models.MockRepos = models.OfficialRepos }()
	session := models.NewMockSession(models.MockDBConfig{})

	chartsImplementation := NewCachedCharts(session)
	freshness := time.Duration(3600) * time.Second
	job := NewRefreshChartsData(chartsImplementation, freshness, "test-run", true)
	assert.Equal(t, job.FirstRun(), true, "First run")
	err := job.Do()
	assert.ExistsErr(t, err, "Error executing refresh")
}
