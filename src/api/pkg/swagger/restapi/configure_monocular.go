package restapi

import (
	"net/http"

	errors "github.com/go-swagger/go-swagger/errors"
	httpkit "github.com/go-swagger/go-swagger/httpkit"
	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/helm/monocular/src/api/handlers"
	"github.com/helm/monocular/src/api/pkg/swagger/restapi/operations"
)

// This file is safe to edit. Once it exists it will not be overwritten

func configureFlags(api *operations.MonocularAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.MonocularAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = httpkit.JSONConsumer()

	api.JSONProducer = httpkit.JSONProducer()

	api.GetChartHandler = operations.GetChartHandlerFunc(func(params operations.GetChartParams) middleware.Responder {
		return handlers.GetChart(params)
	})
	api.GetAllChartsHandler = operations.GetAllChartsHandlerFunc(func() middleware.Responder {
		return handlers.GetAllCharts()
	})
	api.HealthzHandler = operations.HealthzHandlerFunc(func() middleware.Responder {
		return handlers.Healthz()
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
