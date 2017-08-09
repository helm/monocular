package repos

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
	"github.com/go-openapi/runtime"
	"github.com/kubernetes-helm/monocular/src/api/config"
	"github.com/kubernetes-helm/monocular/src/api/data/cache"
	"github.com/kubernetes-helm/monocular/src/api/data/util"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
	reposapi "github.com/kubernetes-helm/monocular/src/api/swagger/restapi/operations/repositories"
	"github.com/kubernetes-helm/monocular/src/api/testutil"
)

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
	assert.Equal(t, len(config.Repos), len(httpBody.Data), "Returns the enabled repos")
}

func setupTestRepoCache() {
	repos := []models.Repo{
		models.Repo{
			Name: util.StrToPtr("stable"),
			URL:  util.StrToPtr("http://storage.googleapis.com/kubernetes-charts"),
		},
		models.Repo{
			Name: util.StrToPtr("incubator"),
			URL:  util.StrToPtr("http://storage.googleapis.com/kubernetes-charts-incubator"),
		},
	}
	cache.NewCachedRepos(repos)
}

func teardownTestRepoCache() {
	if _, err := cache.Repos.DeleteAll(); err != nil {
		log.Fatal("could not clear cache")
	}
}
