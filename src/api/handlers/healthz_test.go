package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
)

func TestHealthz(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthz", nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	Healthz(res, req)
	assert.Equal(t, res.Code, http.StatusOK, "expect a 200 response code")
}
