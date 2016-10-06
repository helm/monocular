package mocks

import (
	"errors"
	"fmt"
	"log"

	"github.com/helm/monocular/src/api/data"
	"github.com/helm/monocular/src/api/data/helpers"
	"github.com/helm/monocular/src/api/pkg/swagger/models"
)

// mockCharts fulfills the data.Charts interface
type mockCharts struct{}

// NewMockCharts returns a new mockCharts
func NewMockCharts() data.Charts {
	return &mockCharts{}
}

// ChartFromRepo method for mockCharts
func (g *mockCharts) ChartFromRepo(repo, name string) (*models.ChartVersion, error) {
	chart, err := GetChartFromMockRepo(repo, name)
	if err != nil {
		return nil, errors.New("chart not found")
	}
	return chart, nil
}

// AllFromRepo method for mockCharts
func (g *mockCharts) AllFromRepo(repo string) ([]*models.ChartVersion, error) {
	charts, err := GetChartsFromMockRepo(repo)
	if err != nil {
		return nil, fmt.Errorf("charts not found for repo %s", repo)
	}
	return charts, nil
}

// All method for mockCharts
func (g *mockCharts) All() ([]*models.Resource, error) {
	charts, err := GetAllChartsFromMockRepos()
	if err != nil {
		return nil, errors.New("unable to load all charts")
	}
	return charts, nil
}

// GetChartFromMockRepo returns a mock "stable/redis" chart resource
func GetChartFromMockRepo(repo, chartName string) (*models.ChartVersion, error) {
	y, err := getMockRepo(repo)
	if err != nil {
		log.Printf("couldn't load mock repo %s!\n", repo)
		return nil, err
	}
	charts, err := helpers.ParseYAMLRepo(y)
	if err != nil {
		log.Printf("couldn't parse mock repo %s!\n", repo)
		return nil, err
	}
	chart, err := helpers.GetLatestChartVersion(charts, chartName)
	if err != nil {
		return nil, err
	}
	return chart, nil
}

// GetAllChartsFromMockRepos returns mock chart resources from all mock repos
func GetAllChartsFromMockRepos() ([]*models.Resource, error) {
	var ret []*models.Resource
	repos := []string{"stable", "incubator"}
	for _, repo := range repos {
		y, err := getMockRepo(repo)
		if err != nil {
			log.Printf("couldn't load mock repo %s!\n", repo)
			return ret, err
		}
		charts, err := helpers.ParseYAMLRepo(y)
		if err != nil {
			log.Printf("couldn't parse mock repo %s!\n", repo)
			return ret, err
		}
		for _, chart := range charts {
			resource := helpers.MakeChartResource(chart, repo)
			ret = append(ret, resource)
		}
	}
	return ret, nil
}

// GetChartsFromMockRepo returns mock chart resources from the passed-in repo
func GetChartsFromMockRepo(repo string) ([]*models.ChartVersion, error) {
	y, err := getMockRepo(repo)
	if err != nil {
		log.Printf("couldn't load mock repo %s!\n", repo)
		return nil, err
	}
	charts, err := helpers.ParseYAMLRepo(y)
	if err != nil {
		log.Printf("couldn't parse mock repo %s!\n", repo)
		return nil, err
	}
	return charts, nil
}

// getMockRepo is a convenience that loads a yaml repo from the filesystem
func getMockRepo(repo string) ([]byte, error) {
	path, err := getTestDataWd()
	if err != nil {
		return nil, err
	}
	path += fmt.Sprintf("repo-%s.yaml", repo)
	y, err := getYAML(path)
	if err != nil {
		log.Printf("couldn't load mock repo %s!\n", path)
		return nil, err
	}
	return y, nil
}
