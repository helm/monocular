package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
)

func TestInClusterGate(t *testing.T) {
	handler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	req, err := http.NewRequest("GET", "/in-cluster", nil)
	assert.NoErr(t, err)
	tests := []struct {
		name    string
		enabled bool
		want    int
	}{
		{"200 status", true, http.StatusOK},
		{"403 status", false, http.StatusForbidden},
	}
	for _, tt := range tests {
		res := httptest.NewRecorder()
		InClusterGate(tt.enabled)(res, req, handler)
		assert.Equal(t, res.Code, tt.want, tt.name)
	}
}
