package repos

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/arschles/assert"
	"github.com/go-openapi/runtime"
	"github.com/kubernetes-helm/monocular/src/api/config"
	"github.com/kubernetes-helm/monocular/src/api/data/pointerto"
	"github.com/kubernetes-helm/monocular/src/api/storage"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
	reposapi "github.com/kubernetes-helm/monocular/src/api/swagger/restapi/operations/repositories"
	"github.com/kubernetes-helm/monocular/src/api/testutil"
)

func TestMain(m *testing.M) {
	flag.Parse()
	storageDrivers := []string{"redis", "mysql"}
	for _, storageDriver := range storageDrivers {
		err := storage.Init(config.StorageConfig{storageDriver, ""})
		if err != nil {
			fmt.Printf("Failed to initialize storage driver: %v\n", err)
			os.Exit(1)
		}
		returnCode := m.Run()
		if returnCode != 0 {
			os.Exit(returnCode)
		}
	}
	os.Exit(0)
}

func TestGetAllRepos200(t *testing.T) {
	setupTestRepoCache()
	defer teardownTestRepoCache()
	w := httptest.NewRecorder()
	params := reposapi.GetAllReposParams{}
	resp := GetRepos(params)
	assert.NotNil(t, resp, "GetRepos response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	var httpBody models.ResourceArrayData
	assert.NoErr(t, testutil.ResourceArrayDataFromJSON(w.Body, &httpBody))
	config, err := config.GetConfig()
	assert.NoErr(t, err)
	assert.Equal(t, len(httpBody.Data), len(config.Repos), "Returns the enabled repos")
}

func TestGetRepo200(t *testing.T) {
	setupTestRepoCache()
	defer teardownTestRepoCache()
	w := httptest.NewRecorder()
	params := reposapi.GetRepoParams{RepoName: "stable"}
	resp := GetRepo(params)
	assert.NotNil(t, resp, "GetRepo response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	var httpBody models.ResourceData
	assert.NoErr(t, testutil.ResourceDataFromJSON(w.Body, &httpBody))
	assert.Equal(t, *httpBody.Data.ID, params.RepoName, "returns the stable repo")
}

func TestGetRepo404(t *testing.T) {
	setupTestRepoCache()
	defer teardownTestRepoCache()
	w := httptest.NewRecorder()
	params := reposapi.GetRepoParams{RepoName: "inexistant"}
	errResp := GetRepo(params)
	assert.NotNil(t, errResp, "GetRepo response")
	errResp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusNotFound, "expect a 404 response code")
	var httpBody models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody))
	testutil.AssertErrBodyData(t, http.StatusNotFound, "repository", httpBody)
}

func TestCreateRepo201(t *testing.T) {
	setupTestRepoCache()
	defer teardownTestRepoCache()
	w := httptest.NewRecorder()
	testRepo := models.Repo{
		Name:   pointerto.String("repoName"),
		URL:    pointerto.String("http://myrepobucket"),
		Source: "http://github.com/my-repo",
	}
	params := reposapi.CreateRepoParams{Data: &testRepo}
	resp := CreateRepo(params, true)
	assert.NotNil(t, resp, "CreateRepo response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusCreated, "expect a 201 response code")
	var httpBody models.ResourceData
	assert.NoErr(t, testutil.ResourceDataFromJSON(w.Body, &httpBody))
	assert.Equal(t, *httpBody.Data.ID, *testRepo.Name, "returns the stable repo")
	_, found, err := storage.Driver.GetRepo(*testRepo.Name)
	assert.NoErr(t, err)
	assert.True(t, found, "Stable repo not found")
}

func TestCreateRepo400(t *testing.T) {
	setupTestRepoCache()
	defer teardownTestRepoCache()
	w := httptest.NewRecorder()
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
		params := reposapi.CreateRepoParams{Data: &tt.repo}
		resp := CreateRepo(params, true)
		assert.NotNil(t, resp, "CreateRepo response")
		resp.WriteResponse(w, runtime.JSONProducer())
		assert.Equal(t, w.Code, http.StatusBadRequest, "expect a 400 response code")
		var httpBody models.Error
		assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody))
		assert.NotNil(t, httpBody.Message, tt.name+" error response")
		assert.Equal(t, *httpBody.Code, int64(http.StatusBadRequest), "response code in HTTP body data")
		assert.True(t, strings.Contains(*httpBody.Message, tt.errorMsg), "error message in HTTP body data")
		_, found, err := storage.Driver.GetRepo(*testRepo.Name)
		assert.NoErr(t, err)
		assert.False(t, found, "Unexpected repo")
	}
}

func TestCreateRepo403(t *testing.T) {
	setupTestRepoCache()
	defer teardownTestRepoCache()
	w := httptest.NewRecorder()
	testRepo := models.Repo{
		Name:   pointerto.String("repoName"),
		URL:    pointerto.String("http://myrepobucket"),
		Source: "http://github.com/my-repo",
	}
	params := reposapi.CreateRepoParams{Data: &testRepo}
	resp := CreateRepo(params, false)
	assert.NotNil(t, resp, "CreateRepo response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusForbidden, "expect a 403 response code")
	var httpBody models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody))
	assert.Equal(t, *httpBody.Code, int64(http.StatusForbidden), "response code in HTTP body data")
	assert.True(t, strings.Contains(*httpBody.Message, "Feature not enabled"), "error message in HTTP body data")
	_, found, err := storage.Driver.GetRepo(*testRepo.Name)
	assert.NoErr(t, err)
	assert.False(t, found, "Unexpected repo")
}

func TestDeleteRepo200(t *testing.T) {
	setupTestRepoCache()
	defer teardownTestRepoCache()
	w := httptest.NewRecorder()
	params := reposapi.DeleteRepoParams{RepoName: "stable"}
	resp := DeleteRepo(params, true)
	assert.NotNil(t, resp, "DeleteRepo response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusOK, "expect a 200 response code")
	var httpBody models.ResourceData
	assert.NoErr(t, testutil.ResourceDataFromJSON(w.Body, &httpBody))
	assert.Nil(t, httpBody.Data.ID, "deleted repo")
	_, found, err := storage.Driver.GetRepo("stable")
	assert.NoErr(t, err)
	assert.False(t, found, "Unexpected repo")
}

func TestDeleteRepo403(t *testing.T) {
	setupTestRepoCache()
	defer teardownTestRepoCache()
	w := httptest.NewRecorder()
	params := reposapi.DeleteRepoParams{RepoName: "stable"}
	resp := DeleteRepo(params, false)
	assert.NotNil(t, resp, "CreateRepo response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusForbidden, "expect a 403 response code")
	var httpBody models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody))
	assert.Equal(t, *httpBody.Code, int64(http.StatusForbidden), "response code in HTTP body data")
	assert.True(t, strings.Contains(*httpBody.Message, "Feature not enabled"), "error message in HTTP body data")
	_, found, err := storage.Driver.GetRepo("stable")
	assert.NoErr(t, err)
	assert.True(t, found, "Did not find \"stable\" repo")
}

func TestDeleteRepo404(t *testing.T) {
	setupTestRepoCache()
	defer teardownTestRepoCache()
	w := httptest.NewRecorder()
	params := reposapi.DeleteRepoParams{RepoName: "inexistant"}
	resp := DeleteRepo(params, true)
	assert.NotNil(t, resp, "DeleteRepo response")
	resp.WriteResponse(w, runtime.JSONProducer())
	assert.Equal(t, w.Code, http.StatusNotFound, "expect a 404 response code")
	var httpBody models.Error
	assert.NoErr(t, testutil.ErrorModelFromJSON(w.Body, &httpBody))
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
	storage.Driver.MergeRepos(repos)
}

func teardownTestRepoCache() {
	if _, err := storage.Driver.DeleteRepos(); err != nil {
		log.Fatal("Could not clear cache ", err)
	}
}
