package data

import (
	"log"
	"testing"

	"github.com/arschles/assert"
	"github.com/kubernetes-helm/monocular/src/api/data/pointerto"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
)

func TestUpdateCache(t *testing.T) {
	testRepo := models.Repo{
		Name:   pointerto.String("repoName"),
		URL:    pointerto.String("http://myrepobucket"),
		Source: "http://github.com/my-repo",
	}
	testRepo2 := models.Repo{
		Name: pointerto.String("repoName2"),
		URL:  pointerto.String("http://myrepobucket2"),
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
			UpdateCache(tt.repos)
			reposCollection, err := GetRepos()
			assert.NotNil(t, reposCollection, "Repos collection created")
			numRepos, err := reposCollection.Count()
			assert.NoErr(t, err)
			assert.Equal(t, numRepos, tt.numRepos, tt.name)
			for _, r := range tt.repos {
				repo := Repo{}
				err := reposCollection.Find(*r.Name, &repo)
				assert.NoErr(t, err)
				assert.Equal(t, *repo.Name, *r.Name, tt.name)
				assert.Equal(t, *repo.URL, *r.URL, tt.name)
				assert.Equal(t, repo.Source, r.Source, tt.name)
			}
		})
	}
}

func TestRepo_ModelId(t *testing.T) {
	tests := []struct {
		name string
		r    *Repo
		want string
	}{
		{"stable repo id", &Repo{Name: pointerto.String("stable")}, "stable"},
		{"no id (unexpected)", &Repo{}, "<nil>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.r.ModelId(), tt.want, tt.name)
		})
	}
}

func TestRepo_SetModelId(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		r    *Repo
		args args
	}{
		{"stable repo id", &Repo{}, args{"stable"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.SetModelId(tt.args.name)
			assert.Equal(t, *tt.r.Name, tt.args.name, tt.name)
		})
	}
}

func teardownTestRepoCache() {
	reposCollection, err := GetRepos()
	if err != nil {
		log.Fatal("could not get Repos collection ", err)
	}
	_, err = reposCollection.DeleteAll()
	if err != nil {
		log.Fatal("could not clear cache ", err)
	}
}
