package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
	"github.com/go-openapi/runtime"
	"github.com/helm/monocular/src/api/data/helpers"
	"github.com/helm/monocular/src/api/mocks"
	"github.com/helm/monocular/src/api/swagger/models"
	"github.com/helm/monocular/src/api/swagger/restapi/operations"
	"github.com/helm/monocular/src/api/testutil"
)

var chartsImplementation = mocks.NewMockCharts()

func TestGetChart200(t *testing.T) {
	chart, err := chartsImplementation.ChartFromRepo(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	w := httptest.NewRecorder()
	params := operations.GetChartParams{
		Repo:      testutil.RepoName,
		ChartName: testutil.ChartName,
	}
	resp := GetChart(params, chartsImplementation)
	assert.NotNil(t, resp, "GetChart response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	httpBody := new(models.ResourceData)
	assert.NoErr(t, testutil.ResourceDataFromJSON(w.Body, httpBody))
	chartResource := helpers.MakeChartResource(chart, testutil.RepoName)
	AssertChartResourceBodyData(t, chartResource, httpBody)
	AssertChartResourceLinksBodyData(t, httpBody)
}

func TestGetChart404(t *testing.T) {
	w := httptest.NewRecorder()
	bogonParams := operations.GetChartParams{
		Repo:      testutil.BogusRepo,
		ChartName: testutil.ChartName,
	}
	errResp := GetChart(bogonParams, chartsImplementation)
	errResp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusNotFound, "expect a 404 response code")
	var httpBody models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody))
	AssertErrBodyData(t, http.StatusNotFound, "chart", httpBody)
}

func TestGetChartVersion200(t *testing.T) {
	chart, err := chartsImplementation.ChartVersionFromRepo(testutil.RepoName, testutil.ChartName, testutil.ChartVersion)
	assert.NoErr(t, err)
	w := httptest.NewRecorder()
	params := operations.GetChartVersionParams{
		Repo:      testutil.RepoName,
		ChartName: testutil.ChartName,
		Version:   testutil.ChartVersion,
	}
	resp := GetChartVersion(params, chartsImplementation)
	assert.NotNil(t, resp, "GetChartVersion response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	httpBody := new(models.ResourceData)
	assert.NoErr(t, testutil.ResourceDataFromJSON(w.Body, httpBody))
	chartResource := helpers.MakeChartResource(chart, testutil.RepoName)
	AssertChartResourceBodyData(t, chartResource, httpBody)
}

func TestGetChartVersion404(t *testing.T) {
	w := httptest.NewRecorder()
	bogonParams := operations.GetChartVersionParams{
		Repo:      testutil.RepoName,
		ChartName: testutil.ChartName,
		Version:   "99.99.99",
	}
	errResp := GetChartVersion(bogonParams, chartsImplementation)
	errResp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusNotFound, "expect a 404 response code")
	var httpBody models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody))
	AssertErrBodyData(t, http.StatusNotFound, "chart", httpBody)
}

func TestGetChartVersions200(t *testing.T) {
	charts, err := chartsImplementation.ChartVersionsFromRepo(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	w := httptest.NewRecorder()
	params := operations.GetChartVersionsParams{
		Repo:      testutil.RepoName,
		ChartName: testutil.ChartName,
	}
	resp := GetChartVersions(params, chartsImplementation)
	assert.NotNil(t, resp, "GetChartVersions response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	var httpBody models.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(w.Body, &httpBody))
	assert.Equal(t, len(charts), len(httpBody.Data), "number of charts returned")
}

func TestGetChartVersions404(t *testing.T) {
	w := httptest.NewRecorder()
	params := operations.GetChartVersionsParams{
		Repo:      testutil.BogusRepo,
		ChartName: testutil.ChartName,
	}
	resp := GetChartVersions(params, chartsImplementation)
	assert.NotNil(t, resp, "GetChartVersions response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusNotFound, "expect a 404 response code")
	var httpBody models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody))
	AssertErrBodyData(t, http.StatusNotFound, "chart", httpBody)
}

func TestGetAllCharts200(t *testing.T) {
	charts, err := chartsImplementation.All()
	assert.NoErr(t, err)
	w := httptest.NewRecorder()
	params := operations.GetAllChartsParams{}
	resp := GetAllCharts(params, chartsImplementation)
	assert.NotNil(t, resp, "GetAllCharts response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	var httpBody models.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(w.Body, &httpBody))
	assert.Equal(t, len(charts), len(httpBody.Data), "number of charts returned")
}

func TestGetChartsInRepo200(t *testing.T) {
	charts, err := chartsImplementation.AllFromRepo(testutil.RepoName)
	assert.NoErr(t, err)
	w := httptest.NewRecorder()
	params := operations.GetChartsInRepoParams{
		Repo: testutil.RepoName,
	}
	resp := GetChartsInRepo(params, chartsImplementation)
	assert.NotNil(t, resp, "GetChartsInRepo response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	var httpBody models.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(w.Body, &httpBody))
	assert.Equal(t, len(charts), len(httpBody.Data), "number of charts returned")
}

func TestGetChartsInRepo404(t *testing.T) {
	w := httptest.NewRecorder()
	params := operations.GetChartsInRepoParams{
		Repo: testutil.BogusRepo,
	}
	resp := GetChartsInRepo(params, chartsImplementation)
	assert.NotNil(t, resp, "GetChartsInRepo response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusNotFound, "expect a 404 response code")
	var httpBody models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody))
	AssertErrBodyData(t, http.StatusNotFound, "charts", httpBody)
}

func TestChartHTTPBody(t *testing.T) {
	w := httptest.NewRecorder()
	chart, err := chartsImplementation.ChartFromRepo(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	chartResource := helpers.MakeChartResource(chart, testutil.RepoName)
	resp := chartHTTPBody(chartResource)
	assert.NotNil(t, resp, "chartHTTPBody response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	httpBody := new(models.ResourceData)
	assert.NoErr(t, testutil.ResourceDataFromJSON(w.Body, httpBody))
	AssertChartResourceBodyData(t, chartResource, httpBody)
}

func TestChartsHTTPBody(t *testing.T) {
	w := httptest.NewRecorder()
	charts, err := chartsImplementation.All()
	assert.NoErr(t, err)
	resp := chartsHTTPBody(charts)
	assert.NotNil(t, resp, "chartHTTPBody response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	var httpBody models.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(w.Body, &httpBody))
	assert.Equal(t, len(charts), len(httpBody.Data), "number of charts returned")
}
