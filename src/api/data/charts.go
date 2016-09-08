package data

import (
	"errors"

	"github.com/helm/monocular/src/api/mocks"
	"github.com/helm/monocular/src/api/pkg/swagger/models"
)

// GetChart gets the chart associated with the passed-in repo+name
func GetChart(repo, name string) (models.Resource, error) {
	// TODO implement actual "get chart" business logic
	chart, err := mocks.GetChartFromMockRepo(repo, name)
	if err != nil {
		return models.Resource{}, errors.New("chart not found")
	}
	return chart, nil
}
