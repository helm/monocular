package charts

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
	"github.com/go-openapi/runtime"
	"github.com/kubernetes-helm/monocular/src/api/data/helpers"
	"github.com/kubernetes-helm/monocular/src/api/handlers"
	"github.com/kubernetes-helm/monocular/src/api/mocks"
	"github.com/kubernetes-helm/monocular/src/api/models"
	swaggermodels "github.com/kubernetes-helm/monocular/src/api/swagger/models"
	chartsapi "github.com/kubernetes-helm/monocular/src/api/swagger/restapi/operations/charts"
	"github.com/kubernetes-helm/monocular/src/api/testutil"
)

var chartsImplementation = mocks.NewMockCharts(mocks.MockedMethods{})
var dbSession = models.NewMockSession(models.MockDBConfig{})
var db, _ = dbSession.DB()
var chartHandlers = NewChartHandlers(dbSession, chartsImplementation)

func TestGetChart200(t *testing.T) {
	chart, err := chartsImplementation.ChartFromRepo(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	req, err := http.NewRequest("GET", "/v1/charts/"+testutil.RepoName+"/"+testutil.ChartName, nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	params := handlers.Params{
		"repo":      testutil.RepoName,
		"chartName": testutil.ChartName,
	}
	chartHandlers.GetChart(res, req, params)
	assert.Equal(t, res.Code, http.StatusOK, "expect a 200 response code")
	httpBody := new(swaggermodels.ResourceData)
	assert.NoErr(t, testutil.ResourceDataFromJSON(res.Body, httpBody))
	chartResource := helpers.MakeChartResource(db, chart)
	testutil.AssertChartResourceBodyData(t, chartResource, httpBody)
}

func TestGetChart404(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/charts/"+testutil.BogusRepo+"/"+testutil.ChartName, nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	bogonParams := handlers.Params{
		"repo":      testutil.BogusRepo,
		"chartName": testutil.ChartName,
	}
	chartHandlers.GetChart(res, req, bogonParams)
	assert.Equal(t, res.Code, http.StatusNotFound, "expect a 404 response code")
	var httpBody swaggermodels.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody))
	testutil.AssertErrBodyData(t, http.StatusNotFound, ChartResourceName, httpBody)
}

func TestGetChartVersion200(t *testing.T) {
	chart, err := chartsImplementation.ChartVersionFromRepo(testutil.RepoName, testutil.ChartName, testutil.ChartVersionString)
	assert.NoErr(t, err)
	req, err := http.NewRequest("GET", "v1/charts/"+testutil.RepoName+"/"+testutil.ChartName+"/versions/"+testutil.ChartVersionString, nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	params := handlers.Params{
		"repo":      testutil.RepoName,
		"chartName": testutil.ChartName,
		"version":   testutil.ChartVersionString,
	}
	chartHandlers.GetChartVersion(res, req, params)
	assert.Equal(t, res.Code, http.StatusOK, "expect a 200 response code")
	httpBody := new(swaggermodels.ResourceData)
	assert.NoErr(t, testutil.ResourceDataFromJSON(res.Body, httpBody))
	chartResource := helpers.MakeChartVersionResource(db, chart)
	testutil.AssertChartVersionResourceBodyData(t, chartResource, httpBody)
}

func TestGetChartVersion404(t *testing.T) {
	req, err := http.NewRequest("GET", "v1/charts/"+testutil.RepoName+"/"+testutil.ChartName+"/versions/"+testutil.ChartVersionString, nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	bogonParams := handlers.Params{
		"repo":      testutil.RepoName,
		"chartName": testutil.ChartName,
		"version":   "99.99.99",
	}
	chartHandlers.GetChartVersion(res, req, bogonParams)
	assert.Equal(t, res.Code, http.StatusNotFound, "expect a 404 response code")
	var httpBody swaggermodels.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody))
	testutil.AssertErrBodyData(t, http.StatusNotFound, ChartVersionResourceName, httpBody)
}

func TestGetChartVersions200(t *testing.T) {
	charts, err := chartsImplementation.ChartVersionsFromRepo(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	req, err := http.NewRequest("GET", "/v1/charts/"+testutil.RepoName+"/"+testutil.ChartName+"/versions", nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	params := handlers.Params{
		"repo":      testutil.RepoName,
		"chartName": testutil.ChartName,
	}
	chartHandlers.GetChartVersions(res, req, params)
	assert.Equal(t, res.Code, http.StatusOK, "expect a 200 response code")
	var httpBody swaggermodels.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(res.Body, &httpBody))
	assert.Equal(t, len(charts), len(httpBody.Data), "number of charts returned")
}

func TestGetChartVersions404(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/charts/"+testutil.BogusRepo+"/"+testutil.ChartName+"/versions", nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	params := handlers.Params{
		"repo":      testutil.BogusRepo,
		"chartName": testutil.ChartName,
	}
	chartHandlers.GetChartVersions(res, req, params)
	assert.Equal(t, res.Code, http.StatusNotFound, "expect a 404 response code")
	var httpBody swaggermodels.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody))
	testutil.AssertErrBodyData(t, http.StatusNotFound, ChartVersionResourceName, httpBody)
}

func TestGetAllCharts200(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/charts", nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	chartHandlers.GetAllCharts(res, req)
	assert.Equal(t, res.Code, http.StatusOK, "expect a 200 response code")
	var httpBody swaggermodels.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(res.Body, &httpBody))
	charts, err := chartsImplementation.All()
	assert.NoErr(t, err)
	assert.Equal(t, len(helpers.MakeChartResources(db, charts)), len(httpBody.Data), "number of charts returned")
}

func TestGetAllCharts404(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/charts", nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	chImplementation := mocks.NewMockCharts(mocks.MockedMethods{
		All: func() ([]*swaggermodels.ChartPackage, error) {
			var ret []*swaggermodels.ChartPackage
			return ret, errors.New("error getting all charts")
		},
	})
	NewChartHandlers(dbSession, chImplementation).GetAllCharts(res, req)
	assert.Equal(t, res.Code, http.StatusNotFound, "expect a 404 response code")
	var httpBody swaggermodels.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody))
	testutil.AssertErrBodyData(t, http.StatusNotFound, ChartResourceName+"s", httpBody)
}

func TestSearchCharts200(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/charts/search?name=drupal", nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	chartHandlers.SearchCharts(res, req)
	assert.Equal(t, res.Code, http.StatusOK, "expect a 200 response code")
	var httpBody swaggermodels.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(res.Body, &httpBody))
	charts, err := chartsImplementation.Search(chartsapi.SearchChartsParams{Name: "drupal"})
	assert.NoErr(t, err)
	assert.Equal(t, len(helpers.MakeChartResources(db, charts)), len(httpBody.Data), "number of charts returned")
}

func TestSearchCharts404(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/charts/search?name=drupal", nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	chImplementation := mocks.NewMockCharts(mocks.MockedMethods{
		Search: func(params chartsapi.SearchChartsParams) ([]*swaggermodels.ChartPackage, error) {
			var ret []*swaggermodels.ChartPackage
			return ret, errors.New("error searching charts")
		},
	})
	NewChartHandlers(dbSession, chImplementation).SearchCharts(res, req)
	assert.Equal(t, res.Code, http.StatusBadRequest, "expect a 400 response code")
	var httpBody swaggermodels.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody))
	assert.Equal(t, *httpBody.Code, int64(400), "response code in HTTP body data")
	assert.Equal(t, *httpBody.Message, "data.Charts Search() error (error searching charts)", "error message in HTTP body data")
}

func TestGetChartsInRepo200(t *testing.T) {
	charts, err := chartsImplementation.AllFromRepo(testutil.RepoName)
	numCharts := len(helpers.MakeChartResources(db, charts))
	assert.NoErr(t, err)
	req, err := http.NewRequest("GET", "/v1/charts/"+testutil.RepoName, nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	chartHandlers.GetChartsInRepo(res, req, handlers.Params{"repo": testutil.RepoName})
	assert.Equal(t, res.Code, http.StatusOK, "expect a 200 response code")
	var httpBody swaggermodels.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(res.Body, &httpBody))
	assert.Equal(t, numCharts, len(httpBody.Data), "number of charts returned")
}

func TestGetChartsInRepo404(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/charts/"+testutil.BogusRepo, nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	chartHandlers.GetChartsInRepo(res, req, handlers.Params{"repo": testutil.BogusRepo})
	assert.Equal(t, res.Code, http.StatusNotFound, "expect a 404 response code")
	var httpBody swaggermodels.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody))
	testutil.AssertErrBodyData(t, http.StatusNotFound, ChartResourceName+"s", httpBody)
}

func TestChartHTTPBody(t *testing.T) {
	w := httptest.NewRecorder()
	chart, err := chartsImplementation.ChartFromRepo(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	chartResource := helpers.MakeChartResource(db, chart)

	payload := handlers.DataResourceBody(chartResource)
	resp := chartsapi.NewGetChartOK().WithPayload(payload)
	assert.NotNil(t, resp, "chartHTTPBody response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	httpBody := new(swaggermodels.ResourceData)
	assert.NoErr(t, testutil.ResourceDataFromJSON(w.Body, httpBody))
	testutil.AssertChartResourceBodyData(t, chartResource, httpBody)
}

func TestChartsHTTPBody(t *testing.T) {
	w := httptest.NewRecorder()
	charts, err := chartsImplementation.All()
	assert.NoErr(t, err)
	resources := helpers.MakeChartResources(db, charts)
	payload := handlers.DataResourcesBody(resources)
	resp := chartsapi.NewGetAllChartsOK().WithPayload(payload)
	assert.NotNil(t, resp, "chartHTTPBody response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	var httpBody swaggermodels.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(w.Body, &httpBody))
	assert.Equal(t, len(resources), len(httpBody.Data), "number of charts returned")
}

func TestNotFound(t *testing.T) {
	const resource1 = "chart"
	const resource2 = "repo"
	res := httptest.NewRecorder()
	notFound(res, resource1)
	assert.Equal(t, res.Code, http.StatusNotFound, "expect a 404 response code")
	var httpBody1 swaggermodels.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody1))
	testutil.AssertErrBodyData(t, http.StatusNotFound, resource1, httpBody1)
	res2 := httptest.NewRecorder()
	var httpBody2 swaggermodels.Error
	notFound(res2, resource2)
	assert.Equal(t, res2.Code, http.StatusNotFound, "expect a 404 response code")
	assert.NoErr(t, testutil.ErrorModelFromJSON(res2.Body, &httpBody2))
	testutil.AssertErrBodyData(t, http.StatusNotFound, resource2, httpBody2)
}
