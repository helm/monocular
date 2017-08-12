package restapi

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/arschles/assert"
	"github.com/go-openapi/loads"
	"github.com/kubernetes-helm/monocular/src/api/data"
	"github.com/kubernetes-helm/monocular/src/api/data/cache"
	"github.com/kubernetes-helm/monocular/src/api/data/helpers"
	handlerscharts "github.com/kubernetes-helm/monocular/src/api/handlers/charts"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
	"github.com/kubernetes-helm/monocular/src/api/swagger/restapi/operations"
	"github.com/kubernetes-helm/monocular/src/api/testutil"
)

const versionsRouteString = "versions"

var chartsImplementation = getChartsImplementation()

// tests the GET /healthz endpoint
func TestGetHealthz(t *testing.T) {
	defer teardownTestRepoCache()
	srv, err := newServer()
	assert.NoErr(t, err)
	defer srv.Close()
	resp, err := httpGet(srv, "healthz")
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK, "response code")
}

// tests the GET /{:apiVersion}/charts endpoint
func TestGetCharts(t *testing.T) {
	defer teardownTestRepoCache()
	srv, err := newServer()
	assert.NoErr(t, err)
	defer srv.Close()
	chartsImplementation.Refresh()
	charts, err := chartsImplementation.All()
	assert.NoErr(t, err)
	resp, err := httpGet(srv, urlPath("v1", handlerscharts.ChartResourceName+"s"))
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK, "response code")
	var httpBody models.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(resp.Body, &httpBody))
	assert.Equal(t, len(helpers.MakeChartResources(charts)), len(httpBody.Data), "number of charts returned")
}

// tests the GET /{:apiVersion}/charts/{:repo} endpoint 200 response
func TestGetChartsInRepo200(t *testing.T) {
	defer teardownTestRepoCache()
	srv, err := newServer()
	assert.NoErr(t, err)
	defer srv.Close()
	chartsImplementation.Refresh()
	charts, err := chartsImplementation.AllFromRepo(testutil.RepoName)
	numCharts := len(helpers.MakeChartResources(charts))
	assert.NoErr(t, err)
	resp, err := httpGet(srv, urlPath("v1", handlerscharts.ChartResourceName+"s", testutil.RepoName))
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK, "response code")
	var httpBody models.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(resp.Body, &httpBody))
	assert.Equal(t, numCharts, len(httpBody.Data), "number of charts returned")
}

// tests the GET /{:apiVersion}/charts/{:repo} endpoint 404 response
func TestGetChartsInRepo404(t *testing.T) {
	defer teardownTestRepoCache()
	srv, err := newServer()
	assert.NoErr(t, err)
	defer srv.Close()
	resp, err := httpGet(srv, urlPath("v1", handlerscharts.ChartResourceName+"s", testutil.BogusRepo))
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusNotFound, "response code")
	var httpBody models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(resp.Body, &httpBody))
	testutil.AssertErrBodyData(t, http.StatusNotFound, handlerscharts.ChartResourceName+"s", httpBody)
}

// tests the GET /{:apiVersion}/charts/{:repo}/{:chart} endpoint 200 response
func TestGetChartInRepo200(t *testing.T) {
	defer teardownTestRepoCache()
	srv, err := newServer()
	assert.NoErr(t, err)
	defer srv.Close()
	chartsImplementation.Refresh()
	chart, err := chartsImplementation.ChartFromRepo(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	resp, err := httpGet(srv, urlPath("v1", handlerscharts.ChartResourceName+"s", testutil.RepoName, testutil.ChartName))
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK, "response code")
	httpBody := new(models.ResourceData)
	assert.NoErr(t, testutil.ResourceDataFromJSON(resp.Body, httpBody))
	chartResource := helpers.MakeChartResource(chart)
	testutil.AssertChartResourceBodyData(t, chartResource, httpBody)
}

// tests the GET /{:apiVersion}/charts/{:repo}/{:chart} endpoint 404 response
func TestGetChartInRepo404(t *testing.T) {
	defer teardownTestRepoCache()
	srv, err := newServer()
	assert.NoErr(t, err)
	defer srv.Close()
	resp, err := httpGet(srv, urlPath("v1", handlerscharts.ChartResourceName+"s", testutil.BogusRepo, testutil.ChartName))
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusNotFound, "response code")
	var httpBody models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(resp.Body, &httpBody))
	testutil.AssertErrBodyData(t, http.StatusNotFound, handlerscharts.ChartResourceName, httpBody)
}

// tests the GET /{:apiVersion}/charts/{:repo}/{:chart}/version endpoint 200 response
func TestGetChartVersion200(t *testing.T) {
	defer teardownTestRepoCache()
	srv, err := newServer()
	assert.NoErr(t, err)
	defer srv.Close()
	chartsImplementation.Refresh()
	chart, err := chartsImplementation.ChartVersionFromRepo(testutil.RepoName, testutil.ChartName, testutil.ChartVersionString)
	assert.NoErr(t, err)
	resp, err := httpGet(srv, urlPath("v1", handlerscharts.ChartResourceName+"s", testutil.RepoName, testutil.ChartName, versionsRouteString, testutil.ChartVersionString))
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK, "response code")
	httpBody := new(models.ResourceData)
	assert.NoErr(t, testutil.ResourceDataFromJSON(resp.Body, httpBody))
	chartResource := helpers.MakeChartVersionResource(chart)
	testutil.AssertChartVersionResourceBodyData(t, chartResource, httpBody)
}

// tests the GET /{:apiVersion}/charts/{:repo}/{:chart}/version endpoint 404 response
func TestGetChartVersion404(t *testing.T) {
	defer teardownTestRepoCache()
	srv, err := newServer()
	assert.NoErr(t, err)
	defer srv.Close()
	resp, err := httpGet(srv, urlPath("v1", handlerscharts.ChartResourceName+"s", testutil.RepoName, testutil.ChartName, versionsRouteString, "99.99.99"))
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusNotFound, "response code")
	var httpBody models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(resp.Body, &httpBody))
	testutil.AssertErrBodyData(t, http.StatusNotFound, handlerscharts.ChartVersionResourceName, httpBody)
}

// tests the GET /{:apiVersion}/charts/{:repo}/{:chart}/version endpoint 200 response
func TestGetChartVersions200(t *testing.T) {
	defer teardownTestRepoCache()
	srv, err := newServer()
	assert.NoErr(t, err)
	defer srv.Close()
	chartsImplementation.Refresh()
	charts, err := chartsImplementation.ChartVersionsFromRepo(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	resp, err := httpGet(srv, urlPath("v1", handlerscharts.ChartResourceName+"s", testutil.RepoName, testutil.ChartName, versionsRouteString))
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK, "response code")
	var httpBody models.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(resp.Body, &httpBody))
	assert.Equal(t, len(charts), len(httpBody.Data), "number of charts returned")
}

// tests the GET /{:apiVersion}/charts/{:repo}/{:chart}/version endpoint 404 response
func TestGetChartVersions404(t *testing.T) {
	defer teardownTestRepoCache()
	srv, err := newServer()
	assert.NoErr(t, err)
	defer srv.Close()
	resp, err := httpGet(srv, urlPath("v1", handlerscharts.ChartResourceName+"s", testutil.BogusRepo, testutil.ChartName, versionsRouteString))
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusNotFound, "response code")
	var httpBody models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(resp.Body, &httpBody))
	testutil.AssertErrBodyData(t, http.StatusNotFound, handlerscharts.ChartVersionResourceName, httpBody)
}

func newServer() (*httptest.Server, error) {
	swaggerSpec, err := loads.Analyzed(SwaggerJSON, "")
	if err != nil {
		return nil, err
	}
	api := operations.NewMonocularAPI(swaggerSpec)
	return httptest.NewServer(configureAPI(api)), nil
}

func urlPath(ver string, remainder ...string) string {
	return fmt.Sprintf("%s/%s", ver, strings.Join(remainder, "/"))
}

func httpGet(s *httptest.Server, route string) (*http.Response, error) {
	return http.Get(s.URL + "/" + route)
}

func getChartsImplementation() data.Charts {
	chartsImplementation := cache.NewCachedCharts()
	return chartsImplementation
}

func teardownTestRepoCache() {
	reposCollection, err := cache.GetRepos()
	if err != nil {
		log.Fatal("could not get Repos collection ", err)
	}
	_, err = reposCollection.DeleteAll()
	if err != nil {
		log.Fatal("could not clear cache ", err)
	}
}
