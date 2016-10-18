package mocks

import (
	"fmt"
	"log"

	"github.com/helm/monocular/src/api/data"
	"github.com/helm/monocular/src/api/data/helpers"
	"github.com/helm/monocular/src/api/swagger/models"
)

// mockCharts fulfills the data.Charts interface
type mockCharts struct{}

// NewMockCharts returns a new mockCharts
func NewMockCharts() data.Charts {
	return &mockCharts{}
}

// ChartFromRepo method for mockCharts
func (g *mockCharts) ChartFromRepo(repo, name string) (*models.ChartVersion, error) {
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
	chart, err := helpers.GetLatestChartVersion(charts, name)
	if err != nil {
		return nil, err
	}
	return chart, nil
}

// ChartVersionsFromRepo is the interface implementation for data.Charts
// It returns the reference to a single versioned chart (the most recently published version)
func (g *mockCharts) ChartVersionsFromRepo(repo, name string) ([]*models.ChartVersion, error) {
	y, err := getMockRepo(repo)
	if err != nil {
		log.Printf("couldn't load mock repo %s!\n", repo)
		return nil, err
	}
	allCharts, err := helpers.ParseYAMLRepo(y)
	if err != nil {
		log.Printf("couldn't parse mock repo %s!\n", repo)
		return nil, err
	}
	charts, err := helpers.GetChartVersions(allCharts, name)
	if err != nil {
		return nil, err
	}
	return charts, nil
}

// AllFromRepo method for mockCharts
func (g *mockCharts) AllFromRepo(repo string) ([]*models.ChartVersion, error) {
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

// All method for mockCharts
func (g *mockCharts) All() ([]*models.Resource, error) {
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

func (g *mockCharts) Refresh() error {
	return nil
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
