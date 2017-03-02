package restapi

import (
	"crypto/tls"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/NYTimes/gziphandler"
	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/helm/monocular/src/api/config"
	"github.com/helm/monocular/src/api/data/cache"
	"github.com/helm/monocular/src/api/data/cache/charthelper"
	"github.com/helm/monocular/src/api/handlers"
	"github.com/helm/monocular/src/api/jobs"
	"github.com/helm/monocular/src/api/swagger/restapi/operations"
	"github.com/rs/cors"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target ../pkg/swagger --name monocular --spec ../swagger-spec/swagger.yml

func configureFlags(api *operations.MonocularAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.MonocularAPI) http.Handler {
	config, err := config.GetConfig()

	if err != nil {
		log.Fatalf("Can not load configuration %v\n", err)
	}
	// configure the api here
	chartsImplementation := cache.NewCachedCharts(config.Repos)
	freshness := time.Duration(3600) * time.Second
	periodicRefresh := cache.NewRefreshChartsData(chartsImplementation, freshness, "refresh-charts")
	toDo := []jobs.Periodic{periodicRefresh}
	jobs.DoPeriodic(toDo)
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
	api.GetChartVersionHandler = operations.GetChartVersionHandlerFunc(func(params operations.GetChartVersionParams) middleware.Responder {
		return handlers.GetChartVersion(params, chartsImplementation)
	})
	api.GetChartVersionsHandler = operations.GetChartVersionsHandlerFunc(func(params operations.GetChartVersionsParams) middleware.Responder {
		return handlers.GetChartVersions(params, chartsImplementation)
	})
	api.GetAllChartsHandler = operations.GetAllChartsHandlerFunc(func(params operations.GetAllChartsParams) middleware.Responder {
		return handlers.GetAllCharts(params, chartsImplementation)
	})
	api.GetChartsInRepoHandler = operations.GetChartsInRepoHandlerFunc(func(params operations.GetChartsInRepoParams) middleware.Responder {
		return handlers.GetChartsInRepo(params, chartsImplementation)
	})
	api.SearchChartsHandler = operations.SearchChartsHandlerFunc(func(params operations.SearchChartsParams) middleware.Responder {
		return handlers.SearchCharts(params, chartsImplementation)
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
	handler = setupStaticFilesMiddleware(handler)
	handler = setupCorsMiddleware(handler)
	return handler
}

// This middleware serves the files existing under cache.DataDirBase
func setupStaticFilesMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Returns static files under /static
		if strings.Index(r.URL.Path, "/assets/") == 0 {
			// 7 days cache
			w.Header().Set("Cache-Control", "public, max-age=604800")
			fs := http.FileServer(http.Dir(charthelper.DataDirBase()))
			fs = http.StripPrefix("/assets/", gziphandler.GzipHandler(fs))
			fs.ServeHTTP(w, r)
		} else {
			// Fallbacks to chained hander
			next.ServeHTTP(w, r)
		}
	})
}

func setupCorsMiddleware(handler http.Handler) http.Handler {
	config, err := config.GetConfig()

	if err != nil {
		log.Fatalf("Can not load configuration %v\n", err)
	}

	c := cors.New(cors.Options{
		AllowedOrigins: config.Cors.AllowedOrigins,
		// They need to be the same than the Access-Control-Request-Headers so it works
		// on pre-flight requests
		AllowedHeaders:   config.Cors.AllowedHeaders,
		AllowCredentials: true,
	})

	// Insert the middleware
	handler = c.Handler(handler)
	return handler
}
