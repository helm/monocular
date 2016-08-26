package handlers

import (
	"github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/helm/monocular/pkg/swagger/restapi/operations"
)

// Healthz is the handler for the /healthz endpoint
func Healthz() middleware.Responder {
	//TODO implement actual health check business logic
	return operations.NewHealthzOK()
}
