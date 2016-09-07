package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/helm/monocular/src/api/pkg/swagger/models"
	"github.com/helm/monocular/src/api/pkg/swagger/restapi/operations"
)

// notFound is a convenience that contains a swagger-friendly 404 given a resource string
func notFound(resource string) middleware.Responder {
	message := fmt.Sprintf("404 %s not found", resource)
	return operations.NewGetChartDefault(http.StatusNotFound).WithPayload(
		&models.Error{Code: http.StatusNotFound, Message: message},
	)
}
