package handlers

import (
	"log"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/helm/monocular/data"
	"github.com/helm/monocular/pkg/swagger/models"
	"github.com/helm/monocular/pkg/swagger/restapi/operations"
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

// chartHTTPBody is a convenience that returns a swagger-friendly HTTP 200 response with chart body data
func chartHTTPBody(chart models.Chart) middleware.Responder {
	return operations.NewGetChartOK().WithPayload(
		operations.GetChartOKBodyBody{
			Data: &chart,
		},
	)
}
