package data

import (
	"errors"
	"fmt"

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

// GetAllCharts gets all charts from all configured repos
func GetAllCharts() ([]*models.Resource, error) {
	// TODO implement actual "get all charts" business logic
	charts, err := mocks.GetAllChartsFromMockRepos()
	if err != nil {
		return nil, errors.New("unable to load all charts")
	}
	return charts, nil
}

// GetChartsInRepo gets all charts from the passed-in repo
func GetChartsInRepo(repo string) ([]*models.Resource, error) {
	// TODO implement actual "get charts in repo" business logic
	charts, err := mocks.GetChartsFromMockRepo(repo)
	if err != nil {
		return nil, fmt.Errorf("charts not found for repo %s", repo)
	}
	return charts, nil
}
