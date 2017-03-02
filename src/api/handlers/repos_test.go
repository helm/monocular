package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
	"github.com/go-openapi/runtime"
	"github.com/helm/monocular/src/api/config"
	"github.com/helm/monocular/src/api/swagger/models"
	"github.com/helm/monocular/src/api/swagger/restapi/operations"
	"github.com/helm/monocular/src/api/testutil"
)

func TestGetAllRepos200(t *testing.T) {
	w := httptest.NewRecorder()
	params := operations.GetAllReposParams{}
	resp := GetRepos(params)
	assert.NotNil(t, resp, "GetRepos response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	var httpBody models.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(w.Body, &httpBody))
	config, err := config.GetConfig()
	assert.NoErr(t, err)
	assert.Equal(t, len(config.Repos), len(httpBody.Data), "Returns the enabled repos")
}
