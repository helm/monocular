package charts

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/kubernetes-helm/monocular/src/api/chartpackagesort"
	"github.com/kubernetes-helm/monocular/src/api/data"
	"github.com/kubernetes-helm/monocular/src/api/data/helpers"
	"github.com/kubernetes-helm/monocular/src/api/data/pointerto"
	"github.com/kubernetes-helm/monocular/src/api/datastore"
	"github.com/kubernetes-helm/monocular/src/api/handlers"
	"github.com/kubernetes-helm/monocular/src/api/handlers/renderer"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
	chartsapi "github.com/kubernetes-helm/monocular/src/api/swagger/restapi/operations/charts"
)

const (
	// ChartResourceName is the resource type string for a chart
	ChartResourceName = "chart"
	// ChartVersionResourceName is the resource type string for a chart version
	ChartVersionResourceName = "chartVersion"
)

// ChartHandlers defines handlers that serve chart data
type ChartHandlers struct {
	dbSession            datastore.Session
	chartsImplementation data.Charts
}

// NewChartHandlers takes a datastore.Session and data.Charts implementation and returns a ChartHandlers struct
func NewChartHandlers(db datastore.Session, ch data.Charts) *ChartHandlers {
	return &ChartHandlers{db, ch}
}

// GetChart is the handler for the /charts/{repo}/{name} endpoint
func (c *ChartHandlers) GetChart(w http.ResponseWriter, req *http.Request, params handlers.Params) {
	chartPackage, err := c.chartsImplementation.ChartFromRepo(params["repo"], params["chartName"])
	if err != nil {
		log.Printf("data.chartsapi.ChartFromRepo(%s, %s) error (%s)", params["repo"], params["chartName"], err)
		notFound(w, ChartResourceName)
		return
	}
	db, closer := c.dbSession.DB()
	defer closer()
	chartResource := helpers.MakeChartResource(db, chartPackage)

	payload := handlers.DataResourceBody(chartResource)
	renderer.Render.JSON(w, http.StatusOK, payload)
}

// GetChartVersion is the handler for the /charts/{repo}/{name}/versions/{version} endpoint
func (c *ChartHandlers) GetChartVersion(w http.ResponseWriter, req *http.Request, params handlers.Params) {
	chartPackage, err := c.chartsImplementation.ChartVersionFromRepo(params["repo"], params["chartName"], params["version"])
	if err != nil {
		log.Printf("data.chartsapi.ChartVersionFromRepo(%s, %s, %s) error (%s)", params["repo"], params["chartName"], params["version"], err)
		notFound(w, ChartVersionResourceName)
		return
	}
	db, closer := c.dbSession.DB()
	defer closer()
	chartVersionResource := helpers.MakeChartVersionResource(db, chartPackage)
	payload := handlers.DataResourceBody(chartVersionResource)
	renderer.Render.JSON(w, http.StatusOK, payload)
}

// GetChartVersions is the handler for the /charts/{repo}/{name}/versions endpoint
func (c *ChartHandlers) GetChartVersions(w http.ResponseWriter, req *http.Request, params handlers.Params) {
	chartPackages, err := c.chartsImplementation.ChartVersionsFromRepo(params["repo"], params["chartName"])
	if err != nil {
		log.Printf("data.chartsapi.ChartVersionsFromRepo(%s, %s) error (%s)", params["repo"], params["chartName"], err)
		notFound(w, ChartVersionResourceName)
		return
	}

	// Sort by semver reverse order
	sort.Sort(sort.Reverse(chartpackagesort.BySemver(chartPackages)))

	db, closer := c.dbSession.DB()
	defer closer()
	chartVersionResources := helpers.MakeChartVersionResources(db, chartPackages)
	payload := handlers.DataResourcesBody(chartVersionResources)
	renderer.Render.JSON(w, http.StatusOK, payload)
}

// GetAllCharts is the handler for the /charts endpoint
func (c *ChartHandlers) GetAllCharts(w http.ResponseWriter, req *http.Request) {
	charts, err := c.chartsImplementation.All()
	if err != nil {
		log.Printf("data.Charts All() error (%s)", err)
		notFound(w, ChartResourceName+"s")
		return
	}

	// For now we only sort by name
	sort.Sort(chartpackagesort.ByName(charts))
	db, closer := c.dbSession.DB()
	defer closer()
	resources := helpers.MakeChartResources(db, charts)
	payload := handlers.DataResourcesBody(resources)
	renderer.Render.JSON(w, http.StatusOK, payload)
}

// GetChartsInRepo is the handler for the /charts/{repo} endpoint
func (c *ChartHandlers) GetChartsInRepo(w http.ResponseWriter, req *http.Request, params handlers.Params) {
	charts, err := c.chartsImplementation.AllFromRepo(params["repo"])
	if err != nil {
		log.Printf("data.Charts AllFromRepo(%s) error (%s)", params["repo"], err)
		notFound(w, ChartResourceName+"s")
		return
	}
	// For now we only sort by name
	sort.Sort(chartpackagesort.ByName(charts))
	db, closer := c.dbSession.DB()
	defer closer()
	chartsResource := helpers.MakeChartResources(db, charts)
	payload := handlers.DataResourcesBody(chartsResource)
	renderer.Render.JSON(w, http.StatusOK, payload)
}

// SearchCharts is the handler for the /charts/search endpoint
func (c *ChartHandlers) SearchCharts(w http.ResponseWriter, req *http.Request) {
	fmt.Println("query", req.URL)
	fmt.Println(req.URL.Query().Get("name"))
	charts, err := c.chartsImplementation.Search(chartsapi.SearchChartsParams{Name: req.URL.Query().Get("name")})
	if err != nil {
		message := fmt.Sprintf("data.Charts Search() error (%s)", err)
		log.Printf(message)
		renderer.Render.JSON(w, http.StatusBadRequest, models.Error{Code: pointerto.Int64(http.StatusBadRequest), Message: &message})
		return
	}
	db, closer := c.dbSession.DB()
	defer closer()
	resources := helpers.MakeChartResources(db, charts)
	payload := handlers.DataResourcesBody(resources)
	renderer.Render.JSON(w, http.StatusOK, payload)
}

func notFound(w http.ResponseWriter, resource string) {
	message := fmt.Sprintf("404 %s not found", resource)
	renderer.Render.JSON(w, http.StatusNotFound, &models.Error{Code: pointerto.Int64(http.StatusNotFound), Message: &message})
}
