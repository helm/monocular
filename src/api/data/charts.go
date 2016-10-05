package data

import (
	"errors"
	"fmt"

	"github.com/helm/monocular/src/api/mocks"
	"github.com/helm/monocular/src/api/pkg/swagger/models"
)

// Charts is an interface for managing chart data
type Charts interface {
	// will have a GetChart method to retrieve a particular chart from a repo
	GetChart(repo, name string) (models.Resource, error)
	// will have a GetAllFromRepo method to retrieve all charts from a repo
	GetAllFromRepo(repo string) ([]*models.Resource, error)
	// will have a GetAll method to retrieve all charts from all repos
	GetAll() ([]*models.Resource, error)
}

// mockCharts fulfills the Charts interface
type mockCharts struct{}

// NewMockCharts returns a new mockCharts
func NewMockCharts() Charts {
	return &mockCharts{}
}

// GetChart method for mockCharts
func (g *mockCharts) GetChart(repo, name string) (models.Resource, error) {
	chart, err := mocks.GetChartFromMockRepo(repo, name)
	if err != nil {
		return models.Resource{}, errors.New("chart not found")
	}
	return chart, nil
}

// GetAllFromRepo method for mockCharts
func (g *mockCharts) GetAllFromRepo(repo string) ([]*models.Resource, error) {
	charts, err := mocks.GetChartsFromMockRepo(repo)
	if err != nil {
		return nil, fmt.Errorf("charts not found for repo %s", repo)
	}
	return charts, nil
}

// GetChart method for mockCharts
func (g *mockCharts) GetAll() ([]*models.Resource, error) {
	charts, err := mocks.GetAllChartsFromMockRepos()
	if err != nil {
		return nil, errors.New("unable to load all charts")
	}
	return charts, nil
}
