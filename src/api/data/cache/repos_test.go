package cache

import (
	"testing"

	"github.com/arschles/assert"

	"github.com/kubernetes-helm/monocular/src/api/data"
	"github.com/kubernetes-helm/monocular/src/api/data/util"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
)

func TestNewCachedRepos(t *testing.T) {
	testRepo := models.Repo{
		Name:   util.StrToPtr("repoName"),
		URL:    util.StrToPtr("http://myrepobucket"),
		Source: "http://github.com/my-repo",
	}
	testRepo2 := models.Repo{
		Name: util.StrToPtr("repoName2"),
		URL:  util.StrToPtr("http://myrepobucket2"),
	}
	tests := []struct {
		name     string
		repos    []models.Repo
		numRepos int
	}{
		{"no repos", []models.Repo{}, 0},
		{"1 repo", []models.Repo{testRepo}, 1},
		{"2 repos", []models.Repo{testRepo, testRepo2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer teardownTestRepoCache()
			NewCachedRepos(tt.repos)
			assert.NotNil(t, Repos, "Repos collection created")
			numRepos, err := Repos.Count()
			assert.NoErr(t, err)
			assert.Equal(t, numRepos, tt.numRepos, tt.name)
			for _, r := range tt.repos {
				repo := data.Repo{}
				err := Repos.Find(*r.Name, &repo)
				assert.NoErr(t, err)
				assert.Equal(t, *repo.Name, *r.Name, tt.name)
				assert.Equal(t, *repo.URL, *r.URL, tt.name)
				assert.Equal(t, repo.Source, r.Source, tt.name)
			}
		})
	}
}
