package mocks

import (
	"fmt"
	"log"
	"strings"

	"github.com/kubernetes-helm/monocular/src/api/data"
	"github.com/kubernetes-helm/monocular/src/api/data/helpers"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
	chartsapi "github.com/kubernetes-helm/monocular/src/api/swagger/restapi/operations/charts"
)

// MockedMethods contains pointers to mocked implementations of methods
type MockedMethods struct {
	All    func() ([]*models.ChartPackage, error)
	Search func(params chartsapi.SearchChartsParams) ([]*models.ChartPackage, error)
}

// mockCharts fulfills the data.Charts interface
type mockCharts struct {
	mockedMethods MockedMethods
}

// NewMockCharts returns a new mockCharts
func NewMockCharts(m MockedMethods) data.Charts {
	return &mockCharts{m}
}

// ChartFromRepo method for mockCharts
func (c *mockCharts) ChartFromRepo(repo, name string) (*models.ChartPackage, error) {
	y, err := getMockRepo(repo)
	if err != nil {
		log.Printf("couldn't load mock repo %s!\n", repo)
		return nil, err
	}
	charts, err := helpers.ParseYAMLRepo(y, repo)
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

// ChartVersionFromRepo is the interface implementation for data.Charts
// It returns the reference to a single versioned chart
func (c *mockCharts) ChartVersionFromRepo(repo, name, version string) (*models.ChartPackage, error) {
	y, err := getMockRepo(repo)
	if err != nil {
		log.Printf("couldn't load mock repo %s!\n", repo)
		return nil, err
	}
	allCharts, err := helpers.ParseYAMLRepo(y, repo)
	if err != nil {
		log.Printf("couldn't parse mock repo %s!\n", repo)
		return nil, err
	}
	chart, err := helpers.GetChartVersion(allCharts, name, version)
	if err != nil {
		return nil, err
	}
	return chart, nil
}

// ChartVersionsFromRepo is the interface implementation for data.Charts
// It returns the reference to a slice of all versions of a particular chart in a repo
func (c *mockCharts) ChartVersionsFromRepo(repo, name string) ([]*models.ChartPackage, error) {
	y, err := getMockRepo(repo)
	if err != nil {
		log.Printf("couldn't load mock repo %s!\n", repo)
		return nil, err
	}
	allCharts, err := helpers.ParseYAMLRepo(y, repo)
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
func (c *mockCharts) AllFromRepo(repo string) ([]*models.ChartPackage, error) {
	y, err := getMockRepo(repo)
	if err != nil {
		log.Printf("couldn't load mock repo %s!\n", repo)
		return nil, err
	}
	charts, err := helpers.ParseYAMLRepo(y, repo)
	if err != nil {
		log.Printf("couldn't parse mock repo %s!\n", repo)
		return nil, err
	}
	return charts, nil
}

// All method for mockCharts
func (c *mockCharts) All() ([]*models.ChartPackage, error) {
	if c.mockedMethods.All != nil {
		return c.mockedMethods.All()
	}
	var allCharts []*models.ChartPackage
	repos := []string{"stable", "incubator"}
	for _, repo := range repos {
		y, err := getMockRepo(repo)
		if err != nil {
			log.Printf("couldn't load mock repo %s!\n", repo)
			return nil, err
		}
		charts, err := helpers.ParseYAMLRepo(y, repo)
		if err != nil {
			log.Printf("couldn't parse mock repo %s!\n", repo)
			return nil, err
		}
		var chartPackages []*models.ChartPackage
		for _, chart := range charts {
			chartPackages = append(chartPackages, chart)
		}
		allCharts = append(allCharts, chartPackages...)
	}
	return allCharts, nil
}

func (c *mockCharts) Search(params chartsapi.SearchChartsParams) ([]*models.ChartPackage, error) {
	if c.mockedMethods.Search != nil {
		return c.mockedMethods.Search(params)
	}
	var ret []*models.ChartPackage
	charts, _ := c.All()
	for _, chart := range charts {
		if strings.Contains(*chart.Name, params.Name) {
			ret = append(ret, chart)
		}
	}
	return ret, nil
}

func (c *mockCharts) Refresh() error {
	return nil
}

// getMockRepo is a convenience that loads a yaml repo from the filesystem
func getMockRepo(repo string) ([]byte, error) {
	path, _ := getTestDataWd()
	path += fmt.Sprintf("repo-%s.yaml", repo)
	y, err := getYAML(path)
	if err != nil {
		log.Printf("couldn't load mock repo %s!\n", path)
		return nil, err
	}
	return y, nil
}
