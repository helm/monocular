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
	helmclient "github.com/helm/monocular/src/api/data/helm/client"

	"github.com/helm/monocular/src/api/config"
	"github.com/helm/monocular/src/api/data/cache"
	"github.com/helm/monocular/src/api/data/cache/charthelper"
	"github.com/helm/monocular/src/api/handlers"
	hcharts "github.com/helm/monocular/src/api/handlers/charts"
	hreleases "github.com/helm/monocular/src/api/handlers/releases"
	hrepos "github.com/helm/monocular/src/api/handlers/repos"
	"github.com/helm/monocular/src/api/jobs"
	"github.com/helm/monocular/src/api/swagger/restapi/operations"
	"github.com/helm/monocular/src/api/swagger/restapi/operations/charts"
	"github.com/helm/monocular/src/api/swagger/restapi/operations/releases"
	"github.com/helm/monocular/src/api/swagger/restapi/operations/repositories"
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
	// Run foreground repository refresh
	chartsImplementation.Refresh()
	// Setup background index refreshes
	cacheRefreshInterval := config.CacheRefreshInterval
	if cacheRefreshInterval <= 0 {
		cacheRefreshInterval = 3600
	}
	freshness := time.Duration(cacheRefreshInterval) * time.Second
	periodicRefresh := cache.NewRefreshChartsData(chartsImplementation, freshness, "refresh-charts", false)
	toDo := []jobs.Periodic{periodicRefresh}
	jobs.DoPeriodic(toDo)

	api.ServeError = errors.ServeError
	helmClient := helmclient.NewHelmClient()

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// s.api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	// Releases
	api.ReleasesGetAllReleasesHandler = releases.GetAllReleasesHandlerFunc(func(params releases.GetAllReleasesParams) middleware.Responder {
		return hreleases.GetReleases(helmClient, params, config.ReleasesEnabled)
	})

	api.ReleasesGetReleaseHandler = releases.GetReleaseHandlerFunc(func(params releases.GetReleaseParams) middleware.Responder {
		return hreleases.GetRelease(helmClient, params, config.ReleasesEnabled)
	})

	api.ReleasesCreateReleaseHandler = releases.CreateReleaseHandlerFunc(func(params releases.CreateReleaseParams) middleware.Responder {
		return hreleases.CreateRelease(helmClient, params, chartsImplementation, config.ReleasesEnabled)
	})

	api.ReleasesDeleteReleaseHandler = releases.DeleteReleaseHandlerFunc(func(params releases.DeleteReleaseParams) middleware.Responder {
		return hreleases.DeleteRelease(helmClient, params, config.ReleasesEnabled)
	})

	// Repos
	api.RepositoriesGetAllReposHandler = repositories.GetAllReposHandlerFunc(func(params repositories.GetAllReposParams) middleware.Responder {
		return hrepos.GetRepos(params)
	})

	// Charts
	api.ChartsSearchChartsHandler = charts.SearchChartsHandlerFunc(func(params charts.SearchChartsParams) middleware.Responder {
		return hcharts.SearchCharts(params, chartsImplementation)
	})

	api.ChartsGetChartHandler = charts.GetChartHandlerFunc(func(params charts.GetChartParams) middleware.Responder {
		return hcharts.GetChart(params, chartsImplementation)
	})

	api.ChartsGetChartVersionHandler = charts.GetChartVersionHandlerFunc(func(params charts.GetChartVersionParams) middleware.Responder {
		return hcharts.GetChartVersion(params, chartsImplementation)
	})

	api.ChartsGetChartVersionsHandler = charts.GetChartVersionsHandlerFunc(func(params charts.GetChartVersionsParams) middleware.Responder {
		return hcharts.GetChartVersions(params, chartsImplementation)
	})

	api.ChartsGetAllChartsHandler = charts.GetAllChartsHandlerFunc(func(params charts.GetAllChartsParams) middleware.Responder {
		return hcharts.GetAllCharts(params, chartsImplementation)
	})

	api.ChartsGetChartsInRepoHandler = charts.GetChartsInRepoHandlerFunc(func(params charts.GetChartsInRepoParams) middleware.Responder {
		return hcharts.GetChartsInRepo(params, chartsImplementation)
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
	handler = gziphandler.GzipHandler(handler)
	return handler
}

// This middleware serves the files existing under cache.DataDirBase
func setupStaticFilesMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Returns static files under /static
		if strings.Index(r.URL.Path, "/assets/") == 0 {
			w.Header().Set("Cache-Control", "public, max-age=7776000")
			fs := http.FileServer(http.Dir(charthelper.DataDirBase()))
			fs = http.StripPrefix("/assets/", fs)
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
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	})

	// Insert the middleware
	handler = c.Handler(handler)
	return handler
}
