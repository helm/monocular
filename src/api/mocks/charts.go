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
	y, err := getYAML(getMocksWd() + fmt.Sprintf("repo-%s.yaml", repo))
	if err != nil {
		log.Fatalf("couldn't load mock repo!")
		return ret, err
	}
	charts, err := helpers.ParseYAMLRepo(y)
	if err != nil {
		log.Fatalf("couldn't parse mock repo!")
		return ret, err
	}
	chart, err := helpers.GetLatestChartVersion(charts, chartName)
	if err != nil {
		return ret, err
	}
	ret = helpers.MakeChartResource(chart, repo, chart.Version)
	return ret, nil
}
