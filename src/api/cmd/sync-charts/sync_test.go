package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
)

type badHTTPClient struct{}

func (h *badHTTPClient) Do(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	w.WriteHeader(500)
	return w.Result(), nil
}

type goodHTTPClient struct{}

func (h *goodHTTPClient) Do(req *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	w.WriteHeader(200)
	return w.Result(), nil
}

func Test_syncURLValidity(t *testing.T) {
	tests := []struct {
		name    string
		repoURL string
	}{
		{"invalid URL", "not-a-url"},
		{"invalid URL", "https//google.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sync("test", tt.repoURL)
			assert.ExistsErr(t, err, tt.name)
		})
	}

	t.Run("valid url", func(t *testing.T) {
		netClient = &goodHTTPClient{}
		err := sync("test", "https://my.examplerepo.com")
		assert.NoErr(t, err)
	})
}

func Test_syncIndexRequest(t *testing.T) {
	netClient = &badHTTPClient{}
	err := sync("test", "https://my.examplerepo.com")
	assert.ExistsErr(t, err, "failed request")
}
