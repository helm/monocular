package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
	"github.com/kubernetes-helm/monocular/src/api/datastore"
	"github.com/kubernetes-helm/monocular/src/api/mocks"
	"github.com/kubernetes-helm/monocular/src/api/models"
	swaggermodels "github.com/kubernetes-helm/monocular/src/api/swagger/models"
	"github.com/kubernetes-helm/monocular/src/api/testutil"
)

func TestUserHandlers_StarChart(t *testing.T) {
	userClaims := models.UserClaims{User: &models.User{Name: "Rick Sanchez", Email: "rick@sanchez.com"}}
	tests := []struct {
		name      string
		params    Params
		expectErr bool
		dbErr     bool
		wantCode  int
	}{
		{"valid chart", Params{"repo": testutil.RepoName, "chart": testutil.ChartName}, false, false, http.StatusNoContent},
		{"inexistant chart", Params{"repo": testutil.BogusRepo, "chart": testutil.ChartName}, true, false, http.StatusNotFound},
		{"db error", Params{"repo": testutil.RepoName, "chart": testutil.ChartName}, true, true, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chImplementation := mocks.NewMockCharts(mocks.MockedMethods{})
			s := datastore.NewMockSession(nil, tt.dbErr)
			u := NewUserHandlers(s, chImplementation)
			req, err := http.NewRequest("POST", "/v1/user/starred/stable/drupal", nil)
			ctx := context.WithValue(req.Context(), models.UserKey, userClaims)
			req = req.WithContext(ctx)
			assert.NoErr(t, err)
			w := httptest.NewRecorder()
			u.StarChart(w, req, tt.params)

			assert.Equal(t, w.Code, tt.wantCode, "response code")
			if tt.expectErr {
				var httpBody swaggermodels.Error
				assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody))
				assert.Equal(t, *httpBody.Code, int64(tt.wantCode), "response code in body")
			}
		})
	}
}

func TestUserHandlers_UnstarChart(t *testing.T) {
	userClaims := models.UserClaims{User: &models.User{Name: "Rick Sanchez", Email: "rick@sanchez.com"}}
	tests := []struct {
		name      string
		params    Params
		expectErr bool
		dbErr     bool
		wantCode  int
	}{
		{"valid chart", Params{"repo": testutil.RepoName, "chart": testutil.ChartName}, false, false, http.StatusNoContent},
		{"inexistant chart", Params{"repo": testutil.BogusRepo, "chart": testutil.ChartName}, true, false, http.StatusNotFound},
		{"db error", Params{"repo": testutil.RepoName, "chart": testutil.ChartName}, true, true, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chImplementation := mocks.NewMockCharts(mocks.MockedMethods{})
			s := datastore.NewMockSession(nil, tt.dbErr)
			u := NewUserHandlers(s, chImplementation)
			req, err := http.NewRequest("DELETE", "/v1/user/starred/stable/drupal", nil)
			ctx := context.WithValue(req.Context(), models.UserKey, userClaims)
			req = req.WithContext(ctx)
			assert.NoErr(t, err)
			w := httptest.NewRecorder()
			u.UnstarChart(w, req, tt.params)

			assert.Equal(t, w.Code, tt.wantCode, "response code")
			if tt.expectErr {
				var httpBody swaggermodels.Error
				assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody))
				assert.Equal(t, *httpBody.Code, int64(tt.wantCode), "response code in body")
			}
		})
	}
}
