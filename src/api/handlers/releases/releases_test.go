package releases

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
	"github.com/go-openapi/runtime"
	"github.com/helm/monocular/src/api/data/helpers"
	"github.com/helm/monocular/src/api/mocks"
	"github.com/helm/monocular/src/api/swagger/models"
	releasesapi "github.com/helm/monocular/src/api/swagger/restapi/operations/releases"
	"github.com/helm/monocular/src/api/testutil"
)

var helmClient = mocks.NewMockedClient()
var helmClientBroken = mocks.NewMockedBrokenClient()
var chartsImplementation = mocks.NewMockCharts()

func validParams() releasesapi.CreateReleaseParams {
	charts, _ := chartsImplementation.All()
	firstChart := charts[0]
	chartID := fmt.Sprintf("%s/%s", firstChart.Repo, *firstChart.Name)
	return releasesapi.CreateReleaseParams{
		Data: releasesapi.CreateReleaseBody{
			ChartID:      helpers.StrToPtr(chartID),
			ChartVersion: firstChart.Version,
		},
	}
}

func TestGetReleases200(t *testing.T) {
	w := httptest.NewRecorder()
	params := releasesapi.GetAllReleasesParams{}
	resp := GetReleases(helmClient, params)
	assert.NotNil(t, resp, "GetReleases response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	// TODO, test body
}

func TestGetReleases500(t *testing.T) {
	w := httptest.NewRecorder()
	params := releasesapi.GetAllReleasesParams{}
	resp := GetReleases(helmClientBroken, params)
	assert.NotNil(t, resp, "Create response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusInternalServerError, "expect a 500 response code")
}

func TestCreateRelease201(t *testing.T) {
	w := httptest.NewRecorder()
	resp := CreateRelease(helmClient, validParams(), chartsImplementation)
	assert.NotNil(t, resp, "Create response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusCreated, "expect a 201 response code")
}

func TestCreateRelease400(t *testing.T) {
	w := httptest.NewRecorder()
	// No ChartVersion
	params := releasesapi.CreateReleaseParams{
		Data: releasesapi.CreateReleaseBody{
			ChartID: helpers.StrToPtr("waps"),
		},
	}
	resp := CreateRelease(helmClient, params, chartsImplementation)
	assert.NotNil(t, resp, "Create response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusBadRequest, "expect a 400 response code")

	// No ChartId
	params = releasesapi.CreateReleaseParams{
		Data: releasesapi.CreateReleaseBody{
			ChartVersion: helpers.StrToPtr("waps"),
		},
	}
	resp = CreateRelease(helmClient, params, chartsImplementation)
	assert.NotNil(t, resp, "Create response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusBadRequest, "expect a 400  response code")
	// Invalid ChartId
	params = releasesapi.CreateReleaseParams{
		Data: releasesapi.CreateReleaseBody{
			ChartID:      helpers.StrToPtr("foo"),
			ChartVersion: helpers.StrToPtr("waps"),
		},
	}
	resp = CreateRelease(helmClient, params, chartsImplementation)
	assert.NotNil(t, resp, "Create response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusBadRequest, "expect a 400 response code")

	// Chart not found
	params = releasesapi.CreateReleaseParams{
		Data: releasesapi.CreateReleaseBody{
			ChartID:      helpers.StrToPtr("stable/foo"),
			ChartVersion: helpers.StrToPtr("does not exist"),
		},
	}
	resp = CreateRelease(helmClient, params, chartsImplementation)
	assert.NotNil(t, resp, "Create response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusBadRequest, "expect a 401 response code")
}

func TestCreateRelease500(t *testing.T) {
	w := httptest.NewRecorder()
	resp := CreateRelease(helmClientBroken, validParams(), chartsImplementation)
	assert.NotNil(t, resp, "Create response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusInternalServerError, "expect a 500 response code")
}

func TestDeleteRelease200(t *testing.T) {
	w := httptest.NewRecorder()
	// No ChartVersion
	params := releasesapi.DeleteReleaseParams{ReleaseName: "foo"}
	resp := DeleteRelease(helmClient, params)
	assert.NotNil(t, resp, "Delete response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
}

func TestDeleteRelease400(t *testing.T) {
	w := httptest.NewRecorder()
	// No ChartVersion
	params := releasesapi.DeleteReleaseParams{}
	resp := DeleteRelease(helmClientBroken, params)
	assert.NotNil(t, resp, "Delete response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusBadRequest, "expect a 400 response code")
}

func TestMakeReleaseResource(t *testing.T) {
	res := makeReleaseResource(&mocks.Resource)
	assert.NotNil(t, res, "Has content")
	assert.Equal(t, *res.Type, "release", "type property")
	assert.Equal(t, *res.ID, "my-release-name", "id property")
	assert.Equal(t, *res.Attributes.(*models.Release).Namespace, "my-namespace", "namespace")
	assert.Equal(t, *res.Attributes.(*models.Release).ChartVersion, "1.2.3", "version")
	assert.Equal(t, *res.Attributes.(*models.Release).ChartName, "my-chart", "chart name")
	assert.Equal(t, *res.Attributes.(*models.Release).Status, "200", "Status")
	assert.NotNil(t, res.Attributes.(*models.Release).Updated, "Has updated at timestamp")

	res = makeReleaseResource(nil)
	assert.NotNil(t, res, "Has content")
}

func TestError(t *testing.T) {
	const resource1 = "release"
	w := httptest.NewRecorder()
	resp := error(resource1)
	assert.NotNil(t, resp, "error response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusInternalServerError, "expect a 500 response code")
	var httpBody1 models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody1))
}

func TestBadRequest(t *testing.T) {
	const resource1 = "release"
	w := httptest.NewRecorder()
	resp := badRequestError(resource1)
	assert.NotNil(t, resp, "badRequest response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusBadRequest, "expect a 400 response code")
	var httpBody1 models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody1))
}
