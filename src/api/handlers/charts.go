package handlers

import (
	"log"

	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/helm/monocular/src/api/data"
	"github.com/helm/monocular/src/api/data/helpers"
	"github.com/helm/monocular/src/api/swagger/models"
	"github.com/helm/monocular/src/api/swagger/restapi/operations"
)

const (
	// ChartResourceName is the resource type string for a chart
	ChartResourceName = "chart"
	// ChartVersionResourceName is the resource type string for a chart version
	ChartVersionResourceName = "chartVersion"
)

// GetChart is the handler for the /charts/{repo}/{name} endpoint
func GetChart(params operations.GetChartParams, c data.Charts) middleware.Responder {
	chartPackage, err := c.ChartFromRepo(params.Repo, params.ChartName)
	if err != nil {
		log.Printf("data.Charts.ChartFromRepo(%s, %s) error (%s)", params.Repo, params.ChartName, err)
		return notFound(ChartResourceName)
	}
	chartResource := helpers.MakeChartResource(chartPackage)
	return chartHTTPBody(chartResource)
}

// GetChartVersion is the handler for the /charts/{repo}/{name}/versions/{version} endpoint
func GetChartVersion(params operations.GetChartVersionParams, c data.Charts) middleware.Responder {
	chartPackage, err := c.ChartVersionFromRepo(params.Repo, params.ChartName, params.Version)
	if err != nil {
		log.Printf("data.Charts.ChartVersionFromRepo(%s, %s, %s) error (%s)", params.Repo, params.ChartName, params.Version, err)
		return notFound(ChartVersionResourceName)
	}
	chartVersionResource := helpers.MakeChartVersionResource(chartPackage)
	return chartHTTPBody(chartVersionResource)
}

// GetChartVersions is the handler for the /charts/{repo}/{name}/versions endpoint
func GetChartVersions(params operations.GetChartVersionsParams, c data.Charts) middleware.Responder {
	chartPackages, err := c.ChartVersionsFromRepo(params.Repo, params.ChartName)
	if err != nil {
		log.Printf("data.Charts.ChartVersionsFromRepo(%s, %s) error (%s)", params.Repo, params.ChartName, err)
		return notFound(ChartVersionResourceName)
	}
	chartVersionResources := helpers.MakeChartVersionResources(chartPackages)
	return chartsHTTPBody(chartVersionResources)
}

// GetAllCharts is the handler for the /charts endpoint
func GetAllCharts(params operations.GetAllChartsParams, c data.Charts) middleware.Responder {
	charts, err := c.All()
	if err != nil {
		log.Printf("data.Charts All() error (%s)", err)
		return notFound(ChartResourceName + "s")
	}
	resources := helpers.MakeChartResources(charts)
	return chartsHTTPBody(resources)
}

// GetChartsInRepo is the handler for the /charts/{repo} endpoint
func GetChartsInRepo(params operations.GetChartsInRepoParams, c data.Charts) middleware.Responder {
	charts, err := c.AllFromRepo(params.Repo)
	if err != nil {
		log.Printf("data.Charts AllFromRepo(%s) error (%s)", params.Repo, err)
		return notFound(ChartResourceName + "s")
	}
	chartsResource := helpers.MakeChartResources(charts)
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
