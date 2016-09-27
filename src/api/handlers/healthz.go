package handlers

import (
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/helm/monocular/src/api/pkg/swagger/restapi/operations"
)

// Healthz is the handler for the /healthz endpoint
func Healthz(params operations.HealthzParams) middleware.Responder {
	//TODO implement actual health check business logic
	return operations.NewHealthzOK()
}
