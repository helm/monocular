package releases

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
	"github.com/kubernetes-helm/monocular/src/api/data/pointerto"
	"github.com/kubernetes-helm/monocular/src/api/handlers"
	"github.com/kubernetes-helm/monocular/src/api/mocks"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
	releasesapi "github.com/kubernetes-helm/monocular/src/api/swagger/restapi/operations/releases"
	"github.com/kubernetes-helm/monocular/src/api/testutil"
)

var helmClient = mocks.NewMockedClient()
var helmClientBroken = mocks.NewMockedBrokenClient()
var chartsImplementation = mocks.NewMockCharts(mocks.MockedMethods{})
var releaseHandlers = NewReleaseHandlers(chartsImplementation, helmClient)
var brokenReleaseHandlers = NewReleaseHandlers(chartsImplementation, helmClientBroken)

func validParams() releasesapi.CreateReleaseBody {
	charts, _ := chartsImplementation.All()
	firstChart := charts[0]
	chartID := fmt.Sprintf("%s/%s", firstChart.Repo, *firstChart.Name)
	return releasesapi.CreateReleaseBody{
		ChartID:      pointerto.String(chartID),
		ChartVersion: firstChart.Version,
	}
}

func TestGetReleases200(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/releases", nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	releaseHandlers.GetReleases(res, req)
	assert.Equal(t, res.Code, http.StatusOK, "expect a 200 response code")
}

func TestGetReleases500(t *testing.T) {
	req, err := http.NewRequest("GET", "/v1/releases", nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	brokenReleaseHandlers.GetReleases(res, req)
	assert.Equal(t, res.Code, http.StatusInternalServerError, "expect a 500 response code")
}

func TestCreateRelease201(t *testing.T) {
	jsonParams, err := json.Marshal(validParams())
	assert.NoErr(t, err)
	req, err := http.NewRequest("POST", "/v1/releases", bytes.NewBuffer(jsonParams))
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	releaseHandlers.CreateRelease(res, req)
	assert.Equal(t, res.Code, http.StatusCreated, "expect a 201 response code")
}

func TestCreateReleaseErrors(t *testing.T) {
	tests := []struct {
		name           string
		params         releasesapi.CreateReleaseBody
		expectedStatus int
	}{
		{"no chartVersion", releasesapi.CreateReleaseBody{ChartID: pointerto.String("waps")}, http.StatusBadRequest},
		{"no chartId", releasesapi.CreateReleaseBody{ChartVersion: pointerto.String("waps")}, http.StatusBadRequest},
		{"invalid chartId", releasesapi.CreateReleaseBody{ChartID: pointerto.String("foo"), ChartVersion: pointerto.String("0.1.0")}, http.StatusBadRequest},
		{"non existant chart", releasesapi.CreateReleaseBody{ChartID: pointerto.String("stable/foo"), ChartVersion: pointerto.String("does not exist")}, http.StatusNotFound},
	}

	for _, tt := range tests {
		jsonParams, err := json.Marshal(tt.params)
		assert.NoErr(t, err)
		req, err := http.NewRequest("POST", "/v1/releases", bytes.NewBuffer(jsonParams))
		assert.NoErr(t, err)
		res := httptest.NewRecorder()
		releaseHandlers.CreateRelease(res, req)
		assert.Equal(t, res.Code, tt.expectedStatus,
			fmt.Sprintf("got %d, expected %d for request with %s", res.Code, tt.expectedStatus, tt.name))
	}
}

func TestCreateRelease500(t *testing.T) {
	jsonParams, err := json.Marshal(validParams())
	assert.NoErr(t, err)
	req, err := http.NewRequest("POST", "/v1/releases", bytes.NewBuffer(jsonParams))
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	brokenReleaseHandlers.CreateRelease(res, req)
	assert.Equal(t, res.Code, http.StatusInternalServerError, "expect a 500 response code")
}

func TestDeleteRelease200(t *testing.T) {
	releaseName := "foo"
	req, err := http.NewRequest("DELETE", "/v1/releases/"+releaseName, nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	releaseHandlers.DeleteRelease(res, req, handlers.Params{"releaseName": releaseName})
	assert.Equal(t, res.Code, http.StatusOK, "expect a 200 response code")
}

func TestDeleteRelease400(t *testing.T) {
	releaseName := "foo"
	req, err := http.NewRequest("DELETE", "/v1/releases/"+releaseName, nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	brokenReleaseHandlers.DeleteRelease(res, req, handlers.Params{"releaseName": releaseName})
	assert.Equal(t, res.Code, http.StatusBadRequest, "expect a 400 response code")
}

func TestGetRelease200(t *testing.T) {
	releaseName := "foo"
	req, err := http.NewRequest("GET", "/v1/releases/"+releaseName, nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	releaseHandlers.GetRelease(res, req, handlers.Params{"releaseName": releaseName})
	assert.Equal(t, res.Code, http.StatusOK, "expect a 200 response code")
}

func TestGetRelease500(t *testing.T) {
	releaseName := "foo"
	req, err := http.NewRequest("GET", "/v1/releases/"+releaseName, nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	brokenReleaseHandlers.GetRelease(res, req, handlers.Params{"releaseName": releaseName})
	assert.Equal(t, res.Code, http.StatusInternalServerError, "expect a 500 response code")
}

func TestMakeReleaseResource(t *testing.T) {
	res := makeReleaseResource(&mocks.Resource)
	assert.NotNil(t, res, "Has content")
	assert.Equal(t, *res.Type, "release", "type property")
	assert.Equal(t, *res.ID, "my-release-name", "id property")
	assert.Equal(t, *res.Attributes.(*models.Release).Namespace, "my-namespace", "namespace")
	assert.Equal(t, *res.Attributes.(*models.Release).ChartVersion, "1.2.3", "version")
	assert.Equal(t, *res.Attributes.(*models.Release).ChartName, "my-chart", "chart name")
	assert.Equal(t, *res.Attributes.(*models.Release).ChartIcon, "chart-icon", "chart icon")
	assert.Equal(t, *res.Attributes.(*models.Release).Status, "200", "Status")
	assert.NotNil(t, res.Attributes.(*models.Release).Updated, "Has updated at timestamp")

	res = makeReleaseResource(nil)
	assert.NotNil(t, res, "Has content")
}

func TestMakeReleaseExtendedResource(t *testing.T) {
	res := makeReleaseExtendedResource(&mocks.Resource)
	assert.NotNil(t, res, "Has content")
	assert.Equal(t, *res.Type, "release", "type property")
	assert.Equal(t, *res.ID, "my-release-name", "id property")
	assert.Equal(t, *res.Attributes.(*models.ReleaseExtended).Namespace, "my-namespace", "namespace")
	assert.Equal(t, *res.Attributes.(*models.ReleaseExtended).ChartVersion, "1.2.3", "version")
	assert.Equal(t, *res.Attributes.(*models.ReleaseExtended).ChartName, "my-chart", "chart name")
	assert.Equal(t, *res.Attributes.(*models.ReleaseExtended).ChartIcon, "chart-icon", "chart icon")
	assert.Equal(t, *res.Attributes.(*models.ReleaseExtended).Status, "200", "Status")
	assert.NotNil(t, res.Attributes.(*models.ReleaseExtended).Updated, "Has updated at timestamp")
	assert.Equal(t, *res.Attributes.(*models.ReleaseExtended).Notes, "my-notes", "Notes")
	assert.Equal(t, *res.Attributes.(*models.ReleaseExtended).Resources, "my-resources", "Notes")

	res = makeReleaseExtendedResource(nil)
	assert.NotNil(t, res, "Has content")
}

func TestErrorResponse(t *testing.T) {
	const resource1 = "release"
	res := httptest.NewRecorder()
	errorResponse(res, http.StatusBadRequest, resource1)
	assert.Equal(t, res.Code, http.StatusBadRequest, "expect a 400 response code")
	var httpBody1 models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody1))
}
