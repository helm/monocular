package handlers

import (
	"log"

	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/helm/monocular/src/api/data"
	"github.com/helm/monocular/src/api/data/helpers"
	"github.com/helm/monocular/src/api/swagger/models"
	"github.com/helm/monocular/src/api/swagger/restapi/operations"
)

const chartResourceName = "chart"

// GetChart is the handler for the /charts/{repo}/{name} endpoint
func GetChart(params operations.GetChartParams, c data.Charts) middleware.Responder {
	chart, err := c.ChartFromRepo(params.Repo, params.ChartName)
	if err != nil {
		log.Printf("data.Charts.ChartFromRepo(%s, %s) error (%s)", params.Repo, params.ChartName, err)
		return notFound(chartResourceName)
	}
	chartResource := helpers.MakeChartResource(chart, params.Repo)
	return chartHTTPBody(chartResource)
}

// GetChartVersions is the handler for the /charts/{repo}/{name}/versions endpoint
func GetChartVersions(params operations.GetChartVersionsParams, c data.Charts) middleware.Responder {
	charts, err := c.ChartVersionsFromRepo(params.Repo, params.ChartName)
	if err != nil {
		log.Printf("data.Charts.ChartVersionsFromRepo(%s, %s) error (%s)", params.Repo, params.ChartName, err)
		return notFound(chartResourceName)
	}
	chartsResource := helpers.MakeChartsResource(charts, params.Repo)
	return chartsHTTPBody(chartsResource)
}

// GetAllCharts is the handler for the /charts endpoint
func GetAllCharts(params operations.GetAllChartsParams, c data.Charts) middleware.Responder {
	charts, err := c.All()
	if err != nil {
		log.Printf("data.Charts All() error (%s)", err)
		return notFound(chartResourceName + "s")
	}
	return chartsHTTPBody(charts)
}

// GetChartsInRepo is the handler for the /charts/{repo} endpoint
func GetChartsInRepo(params operations.GetChartsInRepoParams, c data.Charts) middleware.Responder {
	charts, err := c.AllFromRepo(params.Repo)
	if err != nil {
		log.Printf("data.Charts AllFromRepo(%s) error (%s)", params.Repo, err)
		return notFound(chartResourceName + "s")
	}
	chartsResource := helpers.MakeChartsResource(charts, params.Repo)
	return chartsHTTPBody(chartsResource)
}

// chartHTTPBody is a convenience that returns a swagger-friendly HTTP 200 response with chart body data
func chartHTTPBody(chart *models.Resource) middleware.Responder {
	resourceData := models.ResourceData{
		Data: chart,
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
