package repos

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/arschles/assert"
	"github.com/kubernetes-helm/monocular/src/api/config"
	"github.com/kubernetes-helm/monocular/src/api/data"
	"github.com/kubernetes-helm/monocular/src/api/data/pointerto"
	"github.com/kubernetes-helm/monocular/src/api/handlers"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
	"github.com/kubernetes-helm/monocular/src/api/testutil"
)

func TestGetAllRepos200(t *testing.T) {
	setupTestRepoCache()
	defer teardownTestRepoCache()
	req, err := http.NewRequest("GET", "/v1/repos", nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	GetRepos(res, req)
	assert.Equal(t, res.Code, http.StatusOK, "expect a 200 response code")
	var httpBody models.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(res.Body, &httpBody))
	config, err := config.GetConfig()
	assert.NoErr(t, err)
	assert.Equal(t, len(httpBody.Data), len(config.Repos), "Returns the enabled repos")
}

func TestGetRepo200(t *testing.T) {
	setupTestRepoCache()
	defer teardownTestRepoCache()
	req, err := http.NewRequest("GET", "/v1/repos/"+testutil.RepoName, nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	params := handlers.Params{"repo": testutil.RepoName}
	GetRepo(res, req, params)
	assert.Equal(t, res.Code, http.StatusOK, "expect a 200 response code")
	var httpBody models.ResourceData
	assert.NoErr(t, testutil.ResourceDataFromJSON(res.Body, &httpBody))
	assert.Equal(t, *httpBody.Data.ID, testutil.RepoName, "returns the stable repo")
}

func TestGetRepo404(t *testing.T) {
	setupTestRepoCache()
	defer teardownTestRepoCache()
	req, err := http.NewRequest("GET", "/v1/repos/"+testutil.BogusRepo, nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	params := handlers.Params{"repo": testutil.BogusRepo}
	GetRepo(res, req, params)
	assert.Equal(t, res.Code, http.StatusNotFound, "expect a 404 response code")
	var httpBody models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody))
	testutil.AssertErrBodyData(t, http.StatusNotFound, "repository", httpBody)
}

func TestCreateRepo201(t *testing.T) {
	setupTestRepoCache()
	defer teardownTestRepoCache()
	testRepo := models.Repo{
		Name:   pointerto.String("repoName"),
		URL:    pointerto.String("http://myrepobucket"),
		Source: "http://github.com/my-repo",
	}
	jsonParams, err := json.Marshal(testRepo)
	assert.NoErr(t, err)
	req, err := http.NewRequest("POST", "/v1/repos", bytes.NewBuffer(jsonParams))
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	CreateRepo(res, req)
	assert.Equal(t, res.Code, http.StatusCreated, "expect a 201 response code")
	var httpBody models.ResourceData
	assert.NoErr(t, testutil.ResourceDataFromJSON(res.Body, &httpBody))
	assert.Equal(t, *httpBody.Data.ID, *testRepo.Name, "returns the stable repo")
	reposCollection, _ := data.GetRepos()
	assert.NoErr(t, reposCollection.Find(*testRepo.Name, &data.Repo{}))
}

func TestCreateRepo400(t *testing.T) {
	setupTestRepoCache()
	defer teardownTestRepoCache()
	testRepo := models.Repo{
		Name:   pointerto.String("repoName"),
		Source: "http://github.com/my-repo",
	}
	badURL := models.Repo{
		Name:   testRepo.Name,
		URL:    pointerto.String("not-a-valid-url"),
		Source: testRepo.Source,
	}
	tests := []struct {
		name     string
		repo     models.Repo
		errorMsg string
	}{
		{"no url", testRepo, "URL in body is required"},
		{"bad url", badURL, "URL is invalid"},
	}

	for _, tt := range tests {
		jsonParams, err := json.Marshal(tt.repo)
		assert.NoErr(t, err)

		req, err := http.NewRequest("POST", "/v1/repos", bytes.NewBuffer(jsonParams))
		assert.NoErr(t, err)
		res := httptest.NewRecorder()
		CreateRepo(res, req)

		assert.Equal(t, res.Code, http.StatusBadRequest, "expect a 400 response code")
		var httpBody models.Error
		assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody))
		assert.NotNil(t, httpBody.Message, tt.name+" error response")
		assert.Equal(t, *httpBody.Code, int64(http.StatusBadRequest), "response code in HTTP body data")
		assert.True(t, strings.Contains(*httpBody.Message, tt.errorMsg), "error message in HTTP body data")
		reposCollection, _ := data.GetRepos()
		assert.ExistsErr(t, reposCollection.Find(*testRepo.Name, &data.Repo{}), "invalid repo")
	}
}

func TestDeleteRepo200(t *testing.T) {
	setupTestRepoCache()
	defer teardownTestRepoCache()
	req, err := http.NewRequest("DELETE", "/v1/repos/"+testutil.RepoName, nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	params := handlers.Params{"repo": testutil.RepoName}
	DeleteRepo(res, req, params)
	assert.Equal(t, res.Code, http.StatusOK, "expect a 200 response code")
	var httpBody models.ResourceData
	assert.NoErr(t, testutil.ResourceDataFromJSON(res.Body, &httpBody))
	assert.Nil(t, httpBody.Data.ID, "deleted repo")
	reposCollection, _ := data.GetRepos()
	assert.ExistsErr(t, reposCollection.Find(testutil.RepoName, &data.Repo{}), "deleted repo")
}

func TestDeleteRepo404(t *testing.T) {
	setupTestRepoCache()
	defer teardownTestRepoCache()
	req, err := http.NewRequest("DELETE", "/v1/repos/"+testutil.BogusRepo, nil)
	assert.NoErr(t, err)
	res := httptest.NewRecorder()
	params := handlers.Params{"repo": testutil.BogusRepo}
	DeleteRepo(res, req, params)
	assert.Equal(t, res.Code, http.StatusNotFound, "expect a 404 response code")
	var httpBody models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(res.Body, &httpBody))
	testutil.AssertErrBodyData(t, http.StatusNotFound, "repository", httpBody)
}

func setupTestRepoCache() {
	repos := []models.Repo{
		{
			Name: pointerto.String("stable"),
			URL:  pointerto.String("http://storage.googleapis.com/kubernetes-charts"),
		},
		{
			Name: pointerto.String("incubator"),
			URL:  pointerto.String("http://storage.googleapis.com/kubernetes-charts-incubator"),
		},
	}
	data.UpdateCache(repos)
}

func teardownTestRepoCache() {
	reposCollection, err := data.GetRepos()
	if err != nil {
		log.Fatal("could not get Repos collection ", err)
	}
	_, err = reposCollection.DeleteAll()
	if err != nil {
		log.Fatal("could not clear cache ", err)
	}
}
