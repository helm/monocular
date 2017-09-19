package charts

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/arschles/assert"
	"github.com/go-openapi/runtime"
	"github.com/kubernetes-helm/monocular/src/api/config"
	"github.com/kubernetes-helm/monocular/src/api/data/helpers"
	"github.com/kubernetes-helm/monocular/src/api/data/pointerto"
	"github.com/kubernetes-helm/monocular/src/api/handlers"
	"github.com/kubernetes-helm/monocular/src/api/mocks"
	"github.com/kubernetes-helm/monocular/src/api/storage"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
	chartsapi "github.com/kubernetes-helm/monocular/src/api/swagger/restapi/operations/charts"
	"github.com/kubernetes-helm/monocular/src/api/testutil"
)

func TestMain(m *testing.M) {
	flag.Parse()
	storageDrivers := []string{"redis", "mysql"}
	for _, storageDriver := range storageDrivers {
		err := storage.Init(config.StorageConfig{storageDriver, ""})
		if err != nil {
			fmt.Printf("Failed to initialize storage driver: %v\n", err)
			os.Exit(1)
		}
		returnCode := m.Run()
		if returnCode != 0 {
			os.Exit(returnCode)
		}
	}
	os.Exit(0)
}

var chartsImplementation = mocks.NewMockCharts(mocks.MockedMethods{})

func TestGetChart200(t *testing.T) {
	chart, err := chartsImplementation.ChartFromRepo(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	w := httptest.NewRecorder()
	params := chartsapi.GetChartParams{
		Repo:      testutil.RepoName,
		ChartName: testutil.ChartName,
	}
	resp := GetChart(params, chartsImplementation)
	assert.NotNil(t, resp, "GetChart response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	httpBody := new(models.ResourceData)
	assert.NoErr(t, testutil.ResourceDataFromJSON(w.Body, httpBody))
	chartResource := helpers.MakeChartResource(chart)
	testutil.AssertChartResourceBodyData(t, chartResource, httpBody)
}

func TestGetChart404(t *testing.T) {
	w := httptest.NewRecorder()
	bogonParams := chartsapi.GetChartParams{
		Repo:      testutil.BogusRepo,
		ChartName: testutil.ChartName,
	}
	errResp := GetChart(bogonParams, chartsImplementation)
	errResp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusNotFound, "expect a 404 response code")
	var httpBody models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody))
	testutil.AssertErrBodyData(t, http.StatusNotFound, ChartResourceName, httpBody)
}

func TestGetChartVersion200(t *testing.T) {
	chart, err := chartsImplementation.ChartVersionFromRepo(testutil.RepoName, testutil.ChartName, testutil.ChartVersionString)
	assert.NoErr(t, err)
	w := httptest.NewRecorder()
	params := chartsapi.GetChartVersionParams{
		Repo:      testutil.RepoName,
		ChartName: testutil.ChartName,
		Version:   testutil.ChartVersionString,
	}
	resp := GetChartVersion(params, chartsImplementation)
	assert.NotNil(t, resp, "GetChartVersion response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	httpBody := new(models.ResourceData)
	assert.NoErr(t, testutil.ResourceDataFromJSON(w.Body, httpBody))
	chartResource := helpers.MakeChartVersionResource(chart)
	testutil.AssertChartVersionResourceBodyData(t, chartResource, httpBody)
}

func TestGetChartVersion404(t *testing.T) {
	w := httptest.NewRecorder()
	bogonParams := chartsapi.GetChartVersionParams{
		Repo:      testutil.RepoName,
		ChartName: testutil.ChartName,
		Version:   "99.99.99",
	}
	errResp := GetChartVersion(bogonParams, chartsImplementation)
	errResp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusNotFound, "expect a 404 response code")
	var httpBody models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody))
	testutil.AssertErrBodyData(t, http.StatusNotFound, ChartVersionResourceName, httpBody)
}

func TestGetChartVersions200(t *testing.T) {
	charts, err := chartsImplementation.ChartVersionsFromRepo(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	w := httptest.NewRecorder()
	params := chartsapi.GetChartVersionsParams{
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
	params := chartsapi.GetChartVersionsParams{
		Repo:      testutil.BogusRepo,
		ChartName: testutil.ChartName,
	}
	resp := GetChartVersions(params, chartsImplementation)
	assert.NotNil(t, resp, "GetChartVersions response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusNotFound, "expect a 404 response code")
	var httpBody models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody))
	testutil.AssertErrBodyData(t, http.StatusNotFound, ChartVersionResourceName, httpBody)
}

func TestGetAllCharts200(t *testing.T) {
	setupTestRepoCache()
	defer teardownTestRepoCache()
	w := httptest.NewRecorder()
	params := chartsapi.GetAllChartsParams{}
	resp := GetAllCharts(params, chartsImplementation)
	assert.NotNil(t, resp, "GetAllCharts response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	var httpBody models.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(w.Body, &httpBody))
	charts, err := chartsImplementation.All()
	assert.NoErr(t, err)
	assert.Equal(t, len(helpers.MakeChartResources(charts)), len(httpBody.Data), "number of charts returned")
}

func TestGetAllCharts404(t *testing.T) {
	w := httptest.NewRecorder()
	chImplementation := mocks.NewMockCharts(mocks.MockedMethods{
		All: func() ([]*models.ChartPackage, error) {
			var ret []*models.ChartPackage
			return ret, errors.New("error getting all charts")
		},
	})
	params := chartsapi.GetAllChartsParams{}
	resp := GetAllCharts(params, chImplementation)
	assert.NotNil(t, resp, "GetAllCharts response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusNotFound, "expect a 404 response code")
	var httpBody models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody))
	testutil.AssertErrBodyData(t, http.StatusNotFound, ChartResourceName+"s", httpBody)
}

func TestSearchCharts200(t *testing.T) {
	setupTestRepoCache()
	defer teardownTestRepoCache()
	w := httptest.NewRecorder()
	params := chartsapi.SearchChartsParams{
		Name: "drupal",
	}
	resp := SearchCharts(params, chartsImplementation)
	assert.NotNil(t, resp, "SearchCharts response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	var httpBody models.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(w.Body, &httpBody))
	charts, err := chartsImplementation.Search(params)
	assert.NoErr(t, err)
	assert.Equal(t, len(helpers.MakeChartResources(charts)), len(httpBody.Data), "number of charts returned")
}

func TestSearchCharts404(t *testing.T) {
	w := httptest.NewRecorder()
	params := chartsapi.SearchChartsParams{
		Name: "drupal",
	}
	chImplementation := mocks.NewMockCharts(mocks.MockedMethods{
		Search: func(params chartsapi.SearchChartsParams) ([]*models.ChartPackage, error) {
			var ret []*models.ChartPackage
			return ret, errors.New("error searching charts")
		},
	})
	resp := SearchCharts(params, chImplementation)
	assert.NotNil(t, resp, "SearchCharts response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusBadRequest, "expect a 400 response code")
	var httpBody models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody))
	assert.Equal(t, *httpBody.Code, int64(400), "response code in HTTP body data")
	assert.Equal(t, *httpBody.Message, "data.Charts Search() error (error searching charts)", "error message in HTTP body data")
}

func TestGetChartsInRepo200(t *testing.T) {
	setupTestRepoCache()
	defer teardownTestRepoCache()
	charts, err := chartsImplementation.AllFromRepo(testutil.RepoName)
	numCharts := len(helpers.MakeChartResources(charts))
	assert.NoErr(t, err)
	w := httptest.NewRecorder()
	params := chartsapi.GetChartsInRepoParams{
		Repo: testutil.RepoName,
	}
	resp := GetChartsInRepo(params, chartsImplementation)
	assert.NotNil(t, resp, "GetChartsInRepo response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	var httpBody models.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(w.Body, &httpBody))
	assert.Equal(t, numCharts, len(httpBody.Data), "number of charts returned")
}

func TestGetChartsInRepo404(t *testing.T) {
	w := httptest.NewRecorder()
	params := chartsapi.GetChartsInRepoParams{
		Repo: testutil.BogusRepo,
	}
	resp := GetChartsInRepo(params, chartsImplementation)
	assert.NotNil(t, resp, "GetChartsInRepo response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusNotFound, "expect a 404 response code")
	var httpBody models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody))
	testutil.AssertErrBodyData(t, http.StatusNotFound, ChartResourceName+"s", httpBody)
}

func TestChartHTTPBody(t *testing.T) {
	w := httptest.NewRecorder()
	chart, err := chartsImplementation.ChartFromRepo(testutil.RepoName, testutil.ChartName)
	assert.NoErr(t, err)
	chartResource := helpers.MakeChartResource(chart)

	payload := handlers.DataResourceBody(chartResource)
	resp := chartsapi.NewGetChartOK().WithPayload(payload)
	assert.NotNil(t, resp, "chartHTTPBody response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	httpBody := new(models.ResourceData)
	assert.NoErr(t, testutil.ResourceDataFromJSON(w.Body, httpBody))
	testutil.AssertChartResourceBodyData(t, chartResource, httpBody)
}

func TestChartsHTTPBody(t *testing.T) {
	setupTestRepoCache()
	defer teardownTestRepoCache()
	w := httptest.NewRecorder()
	charts, err := chartsImplementation.All()
	assert.NoErr(t, err)
	resources := helpers.MakeChartResources(charts)
	payload := handlers.DataResourcesBody(resources)
	resp := chartsapi.NewGetAllChartsOK().WithPayload(payload)
	assert.NotNil(t, resp, "chartHTTPBody response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	var httpBody models.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(w.Body, &httpBody))
	assert.Equal(t, len(resources), len(httpBody.Data), "number of charts returned")
}

func TestNotFound(t *testing.T) {
	const resource1 = "chart"
	const resource2 = "repo"
	w := httptest.NewRecorder()
	resp := notFound(resource1)
	assert.NotNil(t, resp, "notFound response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusNotFound, "expect a 404 response code")
	var httpBody1 models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody1))
	testutil.AssertErrBodyData(t, http.StatusNotFound, resource1, httpBody1)
	w = httptest.NewRecorder()
	var httpBody2 models.Error
	resp2 := notFound(resource2)
	assert.NotNil(t, resp2, "notFound response")
	resp2.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusNotFound, "expect a 404 response code")
	assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody2))
	testutil.AssertErrBodyData(t, http.StatusNotFound, resource2, httpBody2)
}

func setupTestRepoCache() {
	repos := []models.Repo{
		{
			Name: pointerto.String("stable"),
			URL:  pointerto.String("http://storage.googleapis.com/kubernetes-charts"),
		},
		{
			Name: pointerto.String("incubator"),
			URL:  pointerto.String("http://storage.googleapis.com/kubernetes-charts-incubator"),
		},
	}
	storage.Driver.MergeRepos(repos)
}

func teardownTestRepoCache() {
	if _, err := storage.Driver.DeleteRepos(); err != nil {
		log.Fatal("Could not clear cache ", err)
	}
}
