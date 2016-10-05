package restapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/arschles/assert"
	"github.com/go-openapi/loads"
	"github.com/helm/monocular/src/api/data"
	"github.com/helm/monocular/src/api/pkg/swagger/models"
	"github.com/helm/monocular/src/api/pkg/swagger/restapi/operations"
	"github.com/helm/monocular/src/api/pkg/testutil"
)

var chartsImplementation = data.NewMockCharts()

// tests the GET /healthz endpoint
func TestGetHealthz(t *testing.T) {
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
	srv, err := newServer()
	assert.NoErr(t, err)
	defer srv.Close()
	charts, err := chartsImplementation.GetAll()
	assert.NoErr(t, err)
	resp, err := httpGet(srv, urlPath("v1", "charts"))
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK, "response code")
	var httpBody models.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(resp.Body, &httpBody))
	assert.Equal(t, len(charts), len(httpBody.Data), "number of charts returned")
}

// tests the GET /{:apiVersion}/charts/{:repo} endpoint 200 response
func TestGetChartsInRepo200(t *testing.T) {
	srv, err := newServer()
	assert.NoErr(t, err)
	defer srv.Close()
	charts, err := chartsImplementation.GetAllFromRepo(testutil.RepoName)
	assert.NoErr(t, err)
	resp, err := httpGet(srv, urlPath("v1", "charts", testutil.RepoName))
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK, "response code")
	var httpBody models.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(resp.Body, &httpBody))
	assert.Equal(t, len(charts), len(httpBody.Data), "number of charts returned")
}

// tests the GET /{:apiVersion}/charts/{:repo} endpoint 404 response
func TestGetChartsInRepo404(t *testing.T) {
	srv, err := newServer()
	assert.NoErr(t, err)
	defer srv.Close()
	resp, err := httpGet(srv, urlPath("v1", "charts", testutil.BogusRepo))
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusNotFound, "response code")
	var httpBody models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(resp.Body, &httpBody))
	testutil.AssertErrBodyData(t, http.StatusNotFound, "charts", httpBody)
}

// tests the GET /{:apiVersion}/charts/{:repo}/{:chart} endpoint 200 response
func TestGetChartInRepo200(t *testing.T) {
	srv, err := newServer()
	assert.NoErr(t, err)
	defer srv.Close()
	chart, err := chartsImplementation.GetChart(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	resp, err := httpGet(srv, urlPath("v1", "charts", testutil.RepoName, testutil.ChartName))
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK, "response code")
	var httpBody models.ResourceData
	assert.NoErr(t, testutil.ResourceDataFromJSON(resp.Body, &httpBody))
	testutil.AssertChartResourceBodyData(t, chart, httpBody)
}

// tests the GET /{:apiVersion}/charts/{:repo}/{:chart} endpoint 404 response
func TestGetChartInRepo404(t *testing.T) {
	srv, err := newServer()
	assert.NoErr(t, err)
	defer srv.Close()
	resp, err := httpGet(srv, urlPath("v1", "charts", testutil.BogusRepo, testutil.ChartName))
	assert.NoErr(t, err)
	defer resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusNotFound, "response code")
	var httpBody models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(resp.Body, &httpBody))
	testutil.AssertErrBodyData(t, http.StatusNotFound, "chart", httpBody)
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
