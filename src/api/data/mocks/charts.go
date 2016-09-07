package mocks

import "github.com/helm/monocular/src/api/pkg/swagger/models"

// GetMockRedisChart returns a mock "kubernetes/redis" chart
func GetMockRedisChart() models.Chart {
	return models.Chart{
		Type: "chart",
		ID:   "charts/redis",
		Links: &models.ChartLinks{
			Latest: "https://storage.googleapis.com/kubernetes-charts/redis-2.0.0.tgz",
		},
		Attributes: &models.ChartAttributes{
			Name: "redis",
			Home: "https://github.com/kubernetes/charts/redis",
		},
	}
}
