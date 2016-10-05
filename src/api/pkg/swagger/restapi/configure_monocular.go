package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/helm/monocular/src/api/data"
	"github.com/helm/monocular/src/api/handlers"
	"github.com/helm/monocular/src/api/pkg/swagger/restapi/operations"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target ../pkg/swagger --name monocular --spec ../swagger-spec/swagger.yml

func configureFlags(api *operations.MonocularAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.MonocularAPI) http.Handler {
	// configure the api here
	chartsImplementation := data.NewMockCharts()
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// s.api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.GetChartHandler = operations.GetChartHandlerFunc(func(params operations.GetChartParams) middleware.Responder {
		return handlers.GetChart(params, chartsImplementation)
	})
	api.GetAllChartsHandler = operations.GetAllChartsHandlerFunc(func(params operations.GetAllChartsParams) middleware.Responder {
		return handlers.GetAllCharts(params, chartsImplementation)
	})
	api.GetChartsInRepoHandler = operations.GetChartsInRepoHandlerFunc(func(params operations.GetChartsInRepoParams) middleware.Responder {
		return handlers.GetChartsInRepo(params, chartsImplementation)
	})
	api.HealthzHandler = operations.HealthzHandlerFunc(func(params operations.HealthzParams) middleware.Responder {
		return handlers.Healthz(params)
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
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
