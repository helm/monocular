package main

import (
	"net/http"
	"os"
	"time"

	"github.com/kubernetes-helm/monocular/src/api/data/cache/charthelper"
	"github.com/rs/cors"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/kubernetes-helm/monocular/src/api/config"
	"github.com/kubernetes-helm/monocular/src/api/config/repos"
	"github.com/kubernetes-helm/monocular/src/api/data"
	"github.com/kubernetes-helm/monocular/src/api/data/cache"
	"github.com/kubernetes-helm/monocular/src/api/data/helm/client"
	"github.com/kubernetes-helm/monocular/src/api/handlers"
	"github.com/kubernetes-helm/monocular/src/api/handlers/charts"
	"github.com/kubernetes-helm/monocular/src/api/handlers/releases"
	repoHandlers "github.com/kubernetes-helm/monocular/src/api/handlers/repos"
	"github.com/kubernetes-helm/monocular/src/api/jobs"
	"github.com/kubernetes-helm/monocular/src/api/middleware"
	"github.com/urfave/negroni"
)

func setupRepoCache(repos repos.Repos) {
	// setup initial chart repositories
	if err := data.UpdateCache(repos); err != nil {
		log.WithError(err).Fatalf("Can not configure repository cache")
	}
}

func setupChartsImplementation(conf config.Configuration) data.Charts {
	setupRepoCache(conf.Repos)

	chartsImplementation := cache.NewCachedCharts()
	// Run foreground repository refresh
	chartsImplementation.Refresh()
	// Setup background index refreshes
	cacheRefreshInterval := conf.CacheRefreshInterval
	if cacheRefreshInterval <= 0 {
		cacheRefreshInterval = 3600
	}
	freshness := time.Duration(cacheRefreshInterval) * time.Second
	periodicRefresh := cache.NewRefreshChartsData(chartsImplementation, freshness, "refresh-charts", false)
	toDo := []jobs.Periodic{periodicRefresh}
	jobs.DoPeriodic(toDo)

	return chartsImplementation
}

func setupCors(conf config.Configuration) *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: conf.Cors.AllowedOrigins,
		// They need to be the same than the Access-Control-Request-Headers so it works
		// on pre-flight requests
		AllowedHeaders:   conf.Cors.AllowedHeaders,
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	})
}

func setupRoutes(conf config.Configuration, chartsImplementation data.Charts, helmClient data.Client) http.Handler {
	r := mux.NewRouter()

	// Middleware
	InClusterGate := middleware.InClusterGate(conf.ReleasesEnabled)
	AuthGate := middleware.AuthGate()

	// Healthcheck
	r.Methods("GET").Path("/healthz").HandlerFunc(handlers.Healthz)

	// API v1
	apiv1 := r.PathPrefix("/v1").Subrouter()

	// Chart routes
	chartHandlers := charts.NewChartHandlers(chartsImplementation)
	apiv1.Methods("GET").Path("/charts").HandlerFunc(chartHandlers.GetAllCharts)
	apiv1.Methods("GET").Path("/charts/search").HandlerFunc(chartHandlers.SearchCharts)
	apiv1.Methods("GET").Path("/charts/{repo}").Handler(handlers.WithParams(chartHandlers.GetChartsInRepo))
	apiv1.Methods("GET").Path("/charts/{repo}/{chartName}").Handler(handlers.WithParams(chartHandlers.GetChart))

	// Chart Version routes
	apiv1.Methods("GET").Path("/charts/{repo}/{chartName}/versions").Handler(handlers.WithParams(chartHandlers.GetChartVersions))
	apiv1.Methods("GET").Path("/charts/{repo}/{chartName}/versions/{version}").Handler(handlers.WithParams(chartHandlers.GetChartVersion))

	// Repo routes
	apiv1.Methods("GET").Path("/repos").HandlerFunc(repoHandlers.GetRepos)
	apiv1.Methods("POST").Path("/repos").Handler(negroni.New(
		InClusterGate,
		AuthGate,
		negroni.WrapFunc(repoHandlers.CreateRepo),
	))
	apiv1.Methods("GET").Path("/repos/{repo}").Handler(handlers.WithParams(repoHandlers.GetRepo))
	apiv1.Methods("DELETE").Path("/repos/{repo}").Handler(negroni.New(
		InClusterGate,
		AuthGate,
		negroni.Wrap(handlers.WithParams(repoHandlers.DeleteRepo)),
	))

	// Releases routes
	releaseHandlers := releases.NewReleaseHandlers(chartsImplementation, helmClient)
	releasesRouter := mux.NewRouter()
	apiv1.PathPrefix("/releases").Handler(negroni.New(InClusterGate, AuthGate, negroni.Wrap(releasesRouter)))
	releasesv1 := releasesRouter.PathPrefix("/v1").Subrouter()
	releasesv1.Methods("GET").Path("/releases").HandlerFunc(releaseHandlers.GetReleases)
	releasesv1.Methods("POST").Path("/releases").HandlerFunc(releaseHandlers.CreateRelease)
	releasesv1.Methods("GET").Path("/releases/{releaseName}").Handler(handlers.WithParams(releaseHandlers.GetRelease))
	releasesv1.Methods("DELETE").Path("/releases/{releaseName}").Handler(handlers.WithParams(releaseHandlers.DeleteRelease))

	// Auth routes
	authHandlers, err := handlers.NewAuthHandlers()
	if err != nil {
		log.WithError(err).Warn("authentication is disabled")
	} else {
		r.Methods("GET").Path("/auth").HandlerFunc(authHandlers.InitiateOAuth)
		r.Methods("GET").Path("/auth/github/callback").HandlerFunc(authHandlers.GithubCallback)
		r.Methods("GET").Path("/auth/verify").Handler(negroni.New(AuthGate))
		r.Methods("DELETE").Path("/auth/logout").HandlerFunc(authHandlers.Logout)
	}

	// Serve chart assets
	fs := http.FileServer(http.Dir(charthelper.DataDirBase()))
	fs = http.StripPrefix("/assets/", fs)
	r.PathPrefix("/assets").Handler(negroni.New(
		negroni.WrapFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Cache-Control", "public, max-age=7776000")
		}),
		negroni.Wrap(fs),
	))

	n := negroni.Classic() // Includes some default middlewares
	n.Use(setupCors(conf))
	n.UseHandler(r)
	return n
}

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		log.WithError(err).Fatal("unable to load configuration")
	}

	chartsImplementation := setupChartsImplementation(conf)
	helmClient := client.NewHelmClient()
	router := setupRoutes(conf, chartsImplementation, helmClient)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.WithFields(log.Fields{"addr": addr}).Info("Started Monocular API server")
	http.ListenAndServe(addr, router)
}
