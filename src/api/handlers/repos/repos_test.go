package repos

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
	"github.com/kubernetes-helm/monocular/src/api/handlers"
	"github.com/kubernetes-helm/monocular/src/api/models"
	swaggermodels "github.com/kubernetes-helm/monocular/src/api/swagger/models"
	"github.com/kubernetes-helm/monocular/src/api/testutil"
)

func TestRepoHandlers_ListRepos(t *testing.T) {
	repos := models.OfficialRepos
	tests := []struct {
		name      string
		expectErr bool
		wantCode  int
	}{
		{"repos", false, http.StatusOK},
		{"db error", true, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := models.NewMockSession(models.MockDBConfig{WantErr: tt.expectErr})
			r := NewRepoHandlers(s)
			req, err := http.NewRequest("GET", "/1/repos", nil)
			assert.NoErr(t, err)
			w := httptest.NewRecorder()
			r.ListRepos(w, req)

			assert.Equal(t, w.Code, tt.wantCode, "response code")
			if tt.expectErr {
				var httpBody swaggermodels.Error
				assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody))
				assert.Equal(t, *httpBody.Code, int64(tt.wantCode), "response code in body")
			} else {
				var httpBody swaggermodels.ResourceArrayData
				assert.NoErr(t, testutil.ResourceArrayDataFromJSON(w.Body, &httpBody))
				assert.Equal(t, len(httpBody.Data), len(repos), tt.name)
				assert.Equal(t, *httpBody.Data[0].ID, repos[0].Name, tt.name)
			}
		})
	}
}

func TestRepoHandlers_GetRepo(t *testing.T) {
	tests := []struct {
		name      string
		expectErr bool
		wantCode  int
	}{
		{"repos", false, http.StatusOK},
		{"db error", true, http.StatusNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := models.NewMockSession(models.MockDBConfig{WantErr: tt.expectErr})
			r := NewRepoHandlers(s)
			req, err := http.NewRequest("GET", "/1/repos/{repo}", nil)
			assert.NoErr(t, err)
			w := httptest.NewRecorder()
			params := handlers.Params{"repo": testutil.RepoName}
			r.GetRepo(w, req, params)

			assert.Equal(t, w.Code, tt.wantCode, "response code")
			if tt.expectErr {
				var httpBody swaggermodels.Error
				assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody))
				assert.Equal(t, *httpBody.Code, int64(tt.wantCode), "response code in body")
			} else {
				var httpBody swaggermodels.ResourceData
				assert.NoErr(t, testutil.ResourceDataFromJSON(w.Body, &httpBody))
				assert.Equal(t, *httpBody.Data.ID, testutil.RepoName, tt.name)
			}
		})
	}
}

func TestRepoHandlers_CreateRepo(t *testing.T) {
	tests := []struct {
		name      string
		body      *bytes.Buffer
		expectErr bool
		wantCode  int
	}{
		{"valid repo", testutil.ToJSONBody(t, &models.Repo{Name: "foo", URL: "foo.com", Source: "foo.com"}), false, http.StatusCreated},
		{"garbage", bytes.NewBufferString("{INVALIDJSON"), true, http.StatusBadRequest},
		{"repo without URL", testutil.ToJSONBody(t, &models.Repo{Name: "foo", Source: "foo.com"}), true, http.StatusBadRequest},
		{"repo with invalid URL", testutil.ToJSONBody(t, &models.Repo{Name: "foo", URL: "invalid"}), true, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := models.NewMockSession(models.MockDBConfig{})
			r := NewRepoHandlers(s)
			req, err := http.NewRequest("POST", "/v1/repos", tt.body)
			assert.NoErr(t, err)
			w := httptest.NewRecorder()
			r.CreateRepo(w, req)

			assert.Equal(t, w.Code, tt.wantCode, "response code")
			if tt.expectErr {
				var httpBody swaggermodels.Error
				assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody))
				assert.Equal(t, *httpBody.Code, int64(tt.wantCode), "response code in body")
			} else {
				var httpBody swaggermodels.ResourceData
				assert.NoErr(t, testutil.ResourceDataFromJSON(w.Body, &httpBody))
				assert.Equal(t, *httpBody.Data.ID, "foo", tt.name)
			}
		})
	}
}

func TestRepoHandlers_DeleteRepo(t *testing.T) {
	tests := []struct {
		name      string
		expectErr bool
		wantCode  int
	}{
		{"deleted repo", false, http.StatusOK},
		{"db error", true, http.StatusNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := models.NewMockSession(models.MockDBConfig{WantErr: tt.expectErr})
			r := NewRepoHandlers(s)
			req, err := http.NewRequest("DELETE", "/1/repos/{repo}", nil)
			assert.NoErr(t, err)
			w := httptest.NewRecorder()
			params := handlers.Params{"repo": testutil.RepoName}
			r.DeleteRepo(w, req, params)

			assert.Equal(t, w.Code, tt.wantCode, "response code")
			if tt.expectErr {
				var httpBody swaggermodels.Error
				assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody))
				assert.Equal(t, *httpBody.Code, int64(tt.wantCode), "response code in body")
			} else {
				var httpBody swaggermodels.ResourceData
				assert.NoErr(t, testutil.ResourceDataFromJSON(w.Body, &httpBody))
				assert.Equal(t, *httpBody.Data.ID, testutil.RepoName, tt.name)
			}
		})
	}
}
