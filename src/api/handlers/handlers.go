package handlers

import (
	"fmt"
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/helm/monocular/src/api/data/helpers"
	"github.com/helm/monocular/src/api/swagger/models"
	"github.com/helm/monocular/src/api/swagger/restapi/operations/charts"
)

// notFound is a convenience that contains a swagger-friendly 404 given a resource string
func notFound(resource string) middleware.Responder {
	message := fmt.Sprintf("404 %s not found", resource)
	return charts.NewGetChartDefault(http.StatusNotFound).WithPayload(
		&models.Error{Code: helpers.Int64ToPtr(http.StatusNotFound), Message: &message},
	)
}
