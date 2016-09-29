package handlers

import (
	"log"

	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/helm/monocular/src/api/data"
	"github.com/helm/monocular/src/api/pkg/swagger/models"
	"github.com/helm/monocular/src/api/pkg/swagger/restapi/operations"
)

const chartResourceName = "chart"

// GetChart is the handler for the /charts/{repo}/{name} endpoint
func GetChart(params operations.GetChartParams) middleware.Responder {
	chart, err := data.GetChart(params.Repo, params.ChartName)
	if err != nil {
		log.Printf("data.GetChart error (%s)", err)
		return notFound(chartResourceName)
	}
	return chartHTTPBody(chart)
}

// GetAllCharts is the handler for the /charts endpoint
func GetAllCharts(params operations.GetAllChartsParams) middleware.Responder {
	charts, err := data.GetAllCharts()
	if err != nil {
		log.Printf("data.GetAllCharts error (%s)", err)
		return notFound(chartResourceName + "s")
	}
	return chartsHTTPBody(charts)
}

// GetChartsInRepo is the handler for the /charts/{repo} endpoint
func GetChartsInRepo(params operations.GetChartsInRepoParams) middleware.Responder {
	charts, err := data.GetChartsInRepo(params.Repo)
	if err != nil {
		log.Printf("data.GetAllCharts error (%s)", err)
		return notFound(chartResourceName + "s")
	}
	return chartsHTTPBody(charts)
}

// chartHTTPBody is a convenience that returns a swagger-friendly HTTP 200 response with chart body data
func chartHTTPBody(chart models.Resource) middleware.Responder {
	resourceData := models.ResourceData{
		Data: &chart,
	}
	return operations.NewGetChartOK().WithPayload(&resourceData)
}

// chartsHTTPBody is a convenience that returns a swagger-friendly HTTP 200 response with charts body data
func chartsHTTPBody(charts []*models.Resource) middleware.Responder {
	resourceArrayData := models.ResourceArrayData{
		Data: charts,
	}
	return operations.NewGetAllChartsOK().WithPayload(&resourceArrayData)
}
