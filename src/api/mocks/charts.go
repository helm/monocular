package mocks

import (
	"fmt"
	"log"

	"github.com/helm/monocular/src/api/data/helpers"
	"github.com/helm/monocular/src/api/pkg/swagger/models"
)

// GetChartFromMockRepo returns a mock "stable/redis" chart resource
func GetChartFromMockRepo(repo, chartName string) (models.Resource, error) {
	var ret models.Resource
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
	chart, err := helpers.GetLatestChartVersion(charts, chartName)
	if err != nil {
		return ret, err
	}
	ret = helpers.MakeChartResource(chart, repo, *chart.Version)
	return ret, nil
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
			resource := helpers.MakeChartResource(chart, repo, *chart.Version)
			ret = append(ret, &resource)
		}
	}
	return ret, nil
}

// GetChartsFromMockRepo returns mock chart resources from the passed-in repo
func GetChartsFromMockRepo(repo string) ([]*models.Resource, error) {
	var ret []*models.Resource
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
		resource := helpers.MakeChartResource(chart, repo, *chart.Version)
		ret = append(ret, &resource)
	}
	return ret, nil
}

// getMockRepo is a convenience that loads a yaml repo from the filesystem
func getMockRepo(repo string) ([]byte, error) {
	path, err := getMocksWd()
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
