package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/arschles/assert"
	"github.com/gorilla/sessions"
	"github.com/kubernetes-helm/monocular/src/api/config"
	"github.com/kubernetes-helm/monocular/src/api/data"
	"github.com/kubernetes-helm/monocular/src/api/data/cache"
	"github.com/kubernetes-helm/monocular/src/api/data/helpers"
	"github.com/kubernetes-helm/monocular/src/api/data/pointerto"
	handlerscharts "github.com/kubernetes-helm/monocular/src/api/handlers/charts"
	"github.com/kubernetes-helm/monocular/src/api/mocks"
	"github.com/kubernetes-helm/monocular/src/api/models"
	swaggermodels "github.com/kubernetes-helm/monocular/src/api/swagger/models"
	releasesapi "github.com/kubernetes-helm/monocular/src/api/swagger/restapi/operations/releases"
	"github.com/kubernetes-helm/monocular/src/api/testutil"
)

const versionsRouteString = "versions"

var dbSession = models.NewMockSession(models.MockDBConfig{})
var db, _ = dbSession.DB()
var helmClient = mocks.NewMockedClient()
var helmClientBroken = mocks.NewMockedBrokenClient()
var chartsImplementation = getChartsImplementation()
var sessionStore = sessions.NewCookieStore([]byte("test"))
var conf, _ = config.GetConfig()

// tests the GET /healthz endpoint
func TestGetHealthz(t *testing.T) {
	ts := httptest.NewServer(setupRoutes(conf, chartsImplementation, helmClient, dbSession))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/healthz")
	assert.NoErr(t, err)
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK, "response code")
}

// tests the GET /{:apiVersion}/charts endpoint
func TestGetCharts(t *testing.T) {
	chartsImplementation.Refresh()
	ts := httptest.NewServer(setupRoutes(conf, chartsImplementation, helmClient, dbSession))
	defer ts.Close()
	charts, err := chartsImplementation.All()
	assert.NoErr(t, err)
	res, err := http.Get(urlPath(ts.URL, "v1", handlerscharts.ChartResourceName+"s"))
	assert.NoErr(t, err)
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK, "response code")
	var httpBody swaggermodels.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(res.Body, &httpBody))
	assert.Equal(t, len(helpers.MakeChartResources(db, charts)), len(httpBody.Data), "number of charts returned")
}

// // tests the GET /{:apiVersion}/charts/{:repo} endpoint 200 response
func TestGetChartsInRepo200(t *testing.T) {
	ts := httptest.NewServer(setupRoutes(conf, chartsImplementation, helmClient, dbSession))
	defer ts.Close()
	charts, err := chartsImplementation.AllFromRepo(testutil.RepoName)
	numCharts := len(helpers.MakeChartResources(db, charts))
	assert.NoErr(t, err)
	res, err := http.Get(urlPath(ts.URL, "v1", handlerscharts.ChartResourceName+"s", testutil.RepoName))
	assert.NoErr(t, err)
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK, "response code")
	var httpBody swaggermodels.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(res.Body, &httpBody))
	assert.Equal(t, numCharts, len(httpBody.Data), "number of charts returned")
}

// tests the GET /{:apiVersion}/charts/{:repo} endpoint 404 response
func TestGetChartsInRepo404(t *testing.T) {
	ts := httptest.NewServer(setupRoutes(conf, chartsImplementation, helmClient, dbSession))
	defer ts.Close()
	res, err := http.Get(urlPath(ts.URL, "v1", handlerscharts.ChartResourceName+"s", testutil.BogusRepo))
	assert.NoErr(t, err)
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusNotFound, "response code")
	var httpBody swaggermodels.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody))
	testutil.AssertErrBodyData(t, http.StatusNotFound, handlerscharts.ChartResourceName+"s", httpBody)
}

// tests the GET /{:apiVersion}/charts/{:repo}/{:chart} endpoint 200 response
func TestGetChartInRepo200(t *testing.T) {
	ts := httptest.NewServer(setupRoutes(conf, chartsImplementation, helmClient, dbSession))
	defer ts.Close()
	chart, err := chartsImplementation.ChartFromRepo(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	res, err := http.Get(urlPath(ts.URL, "v1", handlerscharts.ChartResourceName+"s", testutil.RepoName, testutil.ChartName))
	assert.NoErr(t, err)
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK, "response code")
	httpBody := new(swaggermodels.ResourceData)
	assert.NoErr(t, testutil.ResourceDataFromJSON(res.Body, httpBody))
	chartResource := helpers.MakeChartResource(db, chart)
	testutil.AssertChartResourceBodyData(t, chartResource, httpBody)
}

// tests the GET /{:apiVersion}/charts/{:repo}/{:chart} endpoint 404 response
func TestGetChartInRepo404(t *testing.T) {
	ts := httptest.NewServer(setupRoutes(conf, chartsImplementation, helmClient, dbSession))
	defer ts.Close()
	res, err := http.Get(urlPath(ts.URL, "v1", handlerscharts.ChartResourceName+"s", testutil.BogusRepo, testutil.ChartName))
	assert.NoErr(t, err)
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusNotFound, "response code")
	var httpBody swaggermodels.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody))
	testutil.AssertErrBodyData(t, http.StatusNotFound, handlerscharts.ChartResourceName, httpBody)
}

// tests the GET /{:apiVersion}/charts/{:repo}/{:chart}/version/{:version} endpoint 200 response
func TestGetChartVersion200(t *testing.T) {
	ts := httptest.NewServer(setupRoutes(conf, chartsImplementation, helmClient, dbSession))
	defer ts.Close()
	chart, err := chartsImplementation.ChartVersionFromRepo(testutil.RepoName, testutil.ChartName, testutil.ChartVersionString)
	assert.NoErr(t, err)
	res, err := http.Get(urlPath(ts.URL, "v1", handlerscharts.ChartResourceName+"s", testutil.RepoName, testutil.ChartName, versionsRouteString, testutil.ChartVersionString))
	assert.NoErr(t, err)
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK, "response code")
	httpBody := new(swaggermodels.ResourceData)
	assert.NoErr(t, testutil.ResourceDataFromJSON(res.Body, httpBody))
	chartResource := helpers.MakeChartVersionResource(db, chart)
	testutil.AssertChartVersionResourceBodyData(t, chartResource, httpBody)
}

// tests the GET /{:apiVersion}/charts/{:repo}/{:chart}/version/{:version} endpoint 404 response
func TestGetChartVersion404(t *testing.T) {
	ts := httptest.NewServer(setupRoutes(conf, chartsImplementation, helmClient, dbSession))
	defer ts.Close()
	res, err := http.Get(urlPath(ts.URL, "v1", handlerscharts.ChartResourceName+"s", testutil.RepoName, testutil.ChartName, versionsRouteString, "99.99.99"))
	assert.NoErr(t, err)
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusNotFound, "response code")
	var httpBody swaggermodels.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody))
	testutil.AssertErrBodyData(t, http.StatusNotFound, handlerscharts.ChartVersionResourceName, httpBody)
}

// tests the GET /{:apiVersion}/charts/{:repo}/{:chart}/versions endpoint 200 response
func TestGetChartVersions200(t *testing.T) {
	ts := httptest.NewServer(setupRoutes(conf, chartsImplementation, helmClient, dbSession))
	defer ts.Close()
	charts, err := chartsImplementation.ChartVersionsFromRepo(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	res, err := http.Get(urlPath(ts.URL, "v1", handlerscharts.ChartResourceName+"s", testutil.RepoName, testutil.ChartName, versionsRouteString))
	assert.NoErr(t, err)
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK, "response code")
	var httpBody swaggermodels.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(res.Body, &httpBody))
	assert.Equal(t, len(charts), len(httpBody.Data), "number of charts returned")
}

// tests the GET /{:apiVersion}/charts/{:repo}/{:chart}/versions endpoint 404 response
func TestGetChartVersions404(t *testing.T) {
	ts := httptest.NewServer(setupRoutes(conf, chartsImplementation, helmClient, dbSession))
	defer ts.Close()
	res, err := http.Get(urlPath(ts.URL, "v1", handlerscharts.ChartResourceName+"s", testutil.BogusRepo, testutil.ChartName, versionsRouteString))
	assert.NoErr(t, err)
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusNotFound, "response code")
	var httpBody swaggermodels.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody))
	testutil.AssertErrBodyData(t, http.StatusNotFound, handlerscharts.ChartVersionResourceName, httpBody)
}

// tests the GET /{:apiVersion}/repos endpoint 200 response
func TestGetRepos200(t *testing.T) {
	ts := httptest.NewServer(setupRoutes(conf, chartsImplementation, helmClient, dbSession))
	defer ts.Close()
	res, err := http.Get(urlPath(ts.URL, "v1", "repos"))
	assert.NoErr(t, err)
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK, "response code")
	var httpBody swaggermodels.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(res.Body, &httpBody))
	assert.Equal(t, len(httpBody.Data), 2, "number of repos returned")
	assert.Equal(t, *httpBody.Data[0].ID, testutil.RepoName, "repo name is correct")
}

// tests the POST /{:apiVersion}/repos endpoint 201 response
func TestCreateRepo201(t *testing.T) {
	conf.ReleasesEnabled = true
	defer func() { conf.ReleasesEnabled = false }()
	ts := httptest.NewServer(setupRoutes(conf, chartsImplementation, helmClient, dbSession))
	defer ts.Close()
	repoName := "repoName"
	testRepo := swaggermodels.Repo{
		Name:   &repoName,
		URL:    pointerto.String("http://myrepobucket"),
		Source: "http://github.com/my-repo",
	}
	jsonParams, err := json.Marshal(testRepo)
	assert.NoErr(t, err)
	res, err := http.Post(urlPath(ts.URL, "v1", "repos"), "application/json", bytes.NewBuffer(jsonParams))
	assert.NoErr(t, err)
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusCreated, "response code")
	var httpBody swaggermodels.ResourceData
	assert.NoErr(t, testutil.ResourceDataFromJSON(res.Body, &httpBody))
	assert.Equal(t, *httpBody.Data.ID, repoName, "returns the correct repo name")
}

// tests the POST /{:apiVersion}/repos endpoint 403 response
func TestCreateRepo403(t *testing.T) {
	ts := httptest.NewServer(setupRoutes(conf, chartsImplementation, helmClient, dbSession))
	defer ts.Close()
	repoName := "repoName"
	testRepo := swaggermodels.Repo{
		Name:   &repoName,
		URL:    pointerto.String("http://myrepobucket"),
		Source: "http://github.com/my-repo",
	}
	jsonParams, err := json.Marshal(testRepo)
	assert.NoErr(t, err)
	res, err := http.Post(urlPath(ts.URL, "v1", "repos"), "application/json", bytes.NewBuffer(jsonParams))
	assert.NoErr(t, err)
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusForbidden, "response code")
	var httpBody swaggermodels.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody))
	assert.Equal(t, *httpBody.Code, int64(http.StatusForbidden), "response code in HTTP body data")
	assert.Equal(t, *httpBody.Message, "feature not enabled", "error message")
}

// tests the GET /{:apiVersion}/repos/{:repo} endpoint 200 response
func TestGetRepo200(t *testing.T) {
	ts := httptest.NewServer(setupRoutes(conf, chartsImplementation, helmClient, dbSession))
	defer ts.Close()
	res, err := http.Get(urlPath(ts.URL, "v1", "repos", testutil.RepoName))
	assert.NoErr(t, err)
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK, "response code")
	var httpBody swaggermodels.ResourceData
	assert.NoErr(t, testutil.ResourceDataFromJSON(res.Body, &httpBody))
	assert.Equal(t, *httpBody.Data.ID, testutil.RepoName, "repo name is correct")
}

// tests the DELETE /{:apiVersion}/repos/{:repo} endpoint 200 response
func TestDeleteRepo200(t *testing.T) {
	conf.ReleasesEnabled = true
	defer func() { conf.ReleasesEnabled = false }()
	ts := httptest.NewServer(setupRoutes(conf, chartsImplementation, helmClient, dbSession))
	defer ts.Close()
	req, err := http.NewRequest("DELETE", urlPath(ts.URL, "v1", "repos", testutil.RepoName), nil)
	assert.NoErr(t, err)
	client := &http.Client{}
	res, err := client.Do(req)
	assert.NoErr(t, err)
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK, "response code")
	var httpBody swaggermodels.ResourceData
	assert.NoErr(t, testutil.ResourceDataFromJSON(res.Body, &httpBody))
	assert.Equal(t, *httpBody.Data.ID, testutil.RepoName, "deleted repo")
}

// tests the DELETE /{:apiVersion}/repos/{:repo} endpoint 403 response
func TestDeleteRepo403(t *testing.T) {
	ts := httptest.NewServer(setupRoutes(conf, chartsImplementation, helmClient, dbSession))
	defer ts.Close()
	req, err := http.NewRequest("DELETE", urlPath(ts.URL, "v1", "repos", testutil.RepoName), nil)
	assert.NoErr(t, err)
	client := &http.Client{}
	res, err := client.Do(req)
	assert.NoErr(t, err)
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusForbidden, "response code")
	var httpBody swaggermodels.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody))
	assert.Equal(t, *httpBody.Code, int64(http.StatusForbidden), "response code in HTTP body data")
	assert.Equal(t, *httpBody.Message, "feature not enabled", "error message")
}

// tests the GET /{:apiVersion}/releases endpoint 200 response
func TestGetReleases200(t *testing.T) {
	conf.ReleasesEnabled = true
	defer func() { conf.ReleasesEnabled = false }()
	ts := httptest.NewServer(setupRoutes(conf, chartsImplementation, helmClient, dbSession))
	defer ts.Close()
	res, err := http.Get(urlPath(ts.URL, "v1", "releases"))
	assert.NoErr(t, err)
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK, "response code")
}

// tests the GET /{:apiVersion}/releases endpoint 403 response
func TestGetReleases403(t *testing.T) {
	ts := httptest.NewServer(setupRoutes(conf, chartsImplementation, helmClient, dbSession))
	defer ts.Close()
	res, err := http.Get(urlPath(ts.URL, "v1", "releases"))
	assert.NoErr(t, err)
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusForbidden, "response code")
	var httpBody swaggermodels.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody))
	assert.Equal(t, *httpBody.Code, int64(http.StatusForbidden), "response code in HTTP body data")
	assert.Equal(t, *httpBody.Message, "feature not enabled", "error message")
}

// tests the POST /{:apiVersion}/releases endpoint 201 response
func TestCreateRelease201(t *testing.T) {
	conf.ReleasesEnabled = true
	defer func() { conf.ReleasesEnabled = false }()
	ts := httptest.NewServer(setupRoutes(conf, chartsImplementation, helmClient, dbSession))
	defer ts.Close()
	chartID := fmt.Sprintf("%s/%s", testutil.RepoName, testutil.ChartName)
	params := releasesapi.CreateReleaseBody{
		ChartID:      pointerto.String(chartID),
		ChartVersion: pointerto.String(testutil.ChartVersionString),
	}
	jsonParams, err := json.Marshal(params)
	assert.NoErr(t, err)
	res, err := http.Post(urlPath(ts.URL, "v1", "releases"), "application/json", bytes.NewBuffer(jsonParams))
	assert.NoErr(t, err)
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusCreated, "response code")
}

// tests the POST /{:apiVersion}/releases endpoint 403 response
func TestCreateRelease403(t *testing.T) {
	ts := httptest.NewServer(setupRoutes(conf, chartsImplementation, helmClient, dbSession))
	defer ts.Close()
	chartID := fmt.Sprintf("%s/%s", testutil.RepoName, testutil.ChartName)
	params := releasesapi.CreateReleaseBody{
		ChartID:      pointerto.String(chartID),
		ChartVersion: pointerto.String(testutil.ChartVersionString),
	}
	jsonParams, err := json.Marshal(params)
	assert.NoErr(t, err)
	res, err := http.Post(urlPath(ts.URL, "v1", "releases"), "application/json", bytes.NewBuffer(jsonParams))
	assert.NoErr(t, err)
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusForbidden, "response code")
	var httpBody swaggermodels.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody))
	assert.Equal(t, *httpBody.Code, int64(http.StatusForbidden), "response code in HTTP body data")
	assert.Equal(t, *httpBody.Message, "feature not enabled", "error message")
}

// tests the GET /{:apiVersion}/releases/{:releaseName} endpoint 200 response
func TestGetRelease200(t *testing.T) {
	conf.ReleasesEnabled = true
	defer func() { conf.ReleasesEnabled = false }()
	releaseName := "foo"
	ts := httptest.NewServer(setupRoutes(conf, chartsImplementation, helmClient, dbSession))
	defer ts.Close()
	res, err := http.Get(urlPath(ts.URL, "v1", "releases", releaseName))
	assert.NoErr(t, err)
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK, "response code")
}

// tests the GET /{:apiVersion}/releases/{:releaseName} endpoint 403 response
func TestGetRelease403(t *testing.T) {
	releaseName := "foo"
	ts := httptest.NewServer(setupRoutes(conf, chartsImplementation, helmClient, dbSession))
	defer ts.Close()
	res, err := http.Get(urlPath(ts.URL, "v1", "releases", releaseName))
	assert.NoErr(t, err)
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusForbidden, "response code")
	var httpBody swaggermodels.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody))
	assert.Equal(t, *httpBody.Code, int64(http.StatusForbidden), "response code in HTTP body data")
	assert.Equal(t, *httpBody.Message, "feature not enabled", "error message")
}

// tests the DELETE /{:apiVersion}/releases/{:releaseName} endpoint 200 response
func TestDeleteRelease200(t *testing.T) {
	conf.ReleasesEnabled = true
	defer func() { conf.ReleasesEnabled = false }()
	releaseName := "foo"
	ts := httptest.NewServer(setupRoutes(conf, chartsImplementation, helmClient, dbSession))
	defer ts.Close()
	req, err := http.NewRequest("DELETE", urlPath(ts.URL, "v1", "releases", releaseName), nil)
	assert.NoErr(t, err)
	client := http.Client{}
	res, err := client.Do(req)
	assert.NoErr(t, err)
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusOK, "response code")
}

// tests the DELETE /{:apiVersion}/releases/{:releaseName} endpoint 403 response
func TestDeleteRelease403(t *testing.T) {
	releaseName := "foo"
	ts := httptest.NewServer(setupRoutes(conf, chartsImplementation, helmClient, dbSession))
	defer ts.Close()
	req, err := http.NewRequest("DELETE", urlPath(ts.URL, "v1", "releases", releaseName), nil)
	assert.NoErr(t, err)
	client := http.Client{}
	res, err := client.Do(req)
	assert.NoErr(t, err)
	defer res.Body.Close()
	assert.Equal(t, res.StatusCode, http.StatusForbidden, "response code")
	var httpBody swaggermodels.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody))
	assert.Equal(t, *httpBody.Code, int64(http.StatusForbidden), "response code in HTTP body data")
	assert.Equal(t, *httpBody.Message, "feature not enabled", "error message")
}

func TestAuthGatedRoutes(t *testing.T) {
	os.Setenv("MONOCULAR_AUTH_SIGNING_KEY", "secret")
	defer os.Unsetenv("MONOCULAR_AUTH_SIGNING_KEY")
	conf.ReleasesEnabled = true
	defer func() { conf.ReleasesEnabled = false }()
	ts := httptest.NewServer(setupRoutes(conf, chartsImplementation, helmClient, dbSession))
	defer ts.Close()
	tests := []struct {
		method string
		route  string
	}{
		// Repos
		{"POST", "/repos"},
		{"DELETE", "/repos/stable"},
		// Releases
		{"GET", "/releases"},
		{"POST", "/releases"},
		{"GET", "/releases/cold-hands"},
		{"DELETE", "/releases/cold-hands"},
	}
	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, ts.URL+"/v1"+tt.route, nil)
		assert.NoErr(t, err)
		client := http.Client{}
		res, err := client.Do(req)
		assert.NoErr(t, err)
		defer res.Body.Close()
		assert.Equal(t, res.StatusCode, http.StatusUnauthorized, "response code")
	}
}

func urlPath(ver string, remainder ...string) string {
	return fmt.Sprintf("%s/%s", ver, strings.Join(remainder, "/"))
}

func getChartsImplementation() data.Charts {
	chartsImplementation := cache.NewCachedCharts(dbSession)
	chartsImplementation.Refresh()
	return chartsImplementation
}
