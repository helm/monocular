package mocks

import (
	"fmt"

	"github.com/helm/monocular/src/api/pkg/swagger/models"
)

// GetMockRedisChart returns a mock "stable/redis" chart
func GetMockRedisChart() models.Resource {
	data, _ := getYAML(getMocksWd() + "redis-chart-0.1.0.yaml")
	chart, _ := ParseYAMLChartVersion(data)
	return models.Resource{
		Type: "chart",
		ID:   fmt.Sprintf("stable/%s", chart.Name),
		Links: &models.ChartResourceLinks{
			Latest: chart.URL,
			Home:   chart.Home,
		},
		Attributes: &models.ChartResourceAttributes{
			Name:        chart.Name,
			Description: chart.Description,
			Created:     chart.Created,
		},
	}
}
