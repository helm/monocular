package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
	"github.com/go-openapi/runtime"
	"github.com/helm/monocular/src/api/swagger/models"
	"github.com/helm/monocular/src/api/testutil"
)

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
