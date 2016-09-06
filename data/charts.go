package data

import (
	"errors"

	"github.com/helm/monocular/data/mocks"
	"github.com/helm/monocular/pkg/swagger/models"
)

// GetChart gets the chart associated with the passed-in repo+name
func GetChart(repo, name string) (models.Chart, error) {
	// TODO implement actual "get chart" business logic
	return totallyFakeGetChartForDemo(repo, name)
}

func totallyFakeGetChartForDemo(repo, name string) (models.Chart, error) {
	if repo == "kubernetes" && name == "redis" {
		chart := mocks.GetMockRedisChart()
		return chart, nil
	}
	return models.Chart{}, errors.New("chart not found")
}
