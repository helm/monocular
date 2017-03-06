package releases

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
	"github.com/go-openapi/runtime"
	releasesapi "github.com/helm/monocular/src/api/swagger/restapi/operations/releases"
)

func TestGetReleases200(t *testing.T) {
	w := httptest.NewRecorder()
	params := releasesapi.GetAllReleasesParams{}
	resp := GetReleases(params)
	assert.NotNil(t, resp, "GetReleases response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	// TODO, test body
}

//func TestCreateRelease201(t *testing.T) {
//	w := httptest.NewRecorder()
//	params := releasesapi.CreateReleaseParams{}
//	resp := CreateRelease(params)
//	assert.NotNil(t, resp, "Create response")
//	resp.WriteResponse(w, runtime.JSONProducer())
//	assert.Equal(t, w.Code, http.StatusCreated, "expect a 201 response code")
//	// TODO, test body
//}
