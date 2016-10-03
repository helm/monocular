package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
	"github.com/go-openapi/runtime"
	"github.com/helm/monocular/src/api/data"
	"github.com/helm/monocular/src/api/mocks"
	"github.com/helm/monocular/src/api/pkg/swagger/models"
	"github.com/helm/monocular/src/api/pkg/swagger/restapi/operations"
	"github.com/helm/monocular/src/api/pkg/testutil"
)

func TestGetChart200(t *testing.T) {
	chart, err := data.GetChart(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	w := httptest.NewRecorder()
	params := operations.GetChartParams{
		Repo:      testutil.RepoName,
		ChartName: testutil.ChartName,
	}
	resp := GetChart(params)
	assert.NotNil(t, resp, "GetChart response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	var httpBody models.ResourceData
	assert.NoErr(t, testutil.ResourceDataFromJSON(w.Body, &httpBody))
	AssertChartResourceBodyData(t, chart, httpBody)
}

func TestGetChart404(t *testing.T) {
	w := httptest.NewRecorder()
	bogonParams := operations.GetChartParams{
		Repo:      testutil.BogusRepo,
		ChartName: testutil.ChartName,
	}
	errResp := GetChart(bogonParams)
	errResp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusNotFound, "expect a 404 response code")
	var httpBody models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody))
	AssertErrBodyData(t, http.StatusNotFound, "chart", httpBody)
}

func TestGetAllCharts200(t *testing.T) {
	charts, err := data.GetAllCharts()
	assert.NoErr(t, err)
	w := httptest.NewRecorder()
	params := operations.GetAllChartsParams{}
	resp := GetAllCharts(params)
	assert.NotNil(t, resp, "GetAllCharts response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	var httpBody models.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(w.Body, &httpBody))
	assert.Equal(t, len(charts), len(httpBody.Data), "number of charts returned")
}

func TestGetChartsInRepo200(t *testing.T) {
	charts, err := data.GetChartsInRepo(testutil.RepoName)
	assert.NoErr(t, err)
	w := httptest.NewRecorder()
	params := operations.GetChartsInRepoParams{
		Repo: testutil.RepoName,
	}
	resp := GetChartsInRepo(params)
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
		Repo: "bogon",
	}
	resp := GetChartsInRepo(params)
	assert.NotNil(t, resp, "GetChartsInRepo response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusNotFound, "expect a 404 response code")
	var httpBody models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody))
	AssertErrBodyData(t, http.StatusNotFound, "charts", httpBody)
}

func TestChartHTTPBody(t *testing.T) {
	w := httptest.NewRecorder()
	chart, err := mocks.GetChartFromMockRepo(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	resp := chartHTTPBody(chart)
	assert.NotNil(t, resp, "chartHTTPBody response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	var httpBody models.ResourceData
	assert.NoErr(t, testutil.ResourceDataFromJSON(w.Body, &httpBody))
	AssertChartResourceBodyData(t, chart, httpBody)
}

func TestChartsHTTPBody(t *testing.T) {
	w := httptest.NewRecorder()
	charts, err := mocks.GetAllChartsFromMockRepos()
	assert.NoErr(t, err)
	resp := chartsHTTPBody(charts)
	assert.NotNil(t, resp, "chartHTTPBody response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	var httpBody models.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(w.Body, &httpBody))
	assert.Equal(t, len(charts), len(httpBody.Data), "number of charts returned")
}
