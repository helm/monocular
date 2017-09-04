package storage

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/arschles/assert"
	"github.com/kubernetes-helm/monocular/src/api/config"
	"github.com/kubernetes-helm/monocular/src/api/data"
	"github.com/kubernetes-helm/monocular/src/api/data/pointerto"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
)

func TestMain(m *testing.M) {
	flag.Parse()
	storageDrivers := []string{"redis", "mysql"}
	for _, storageDriver := range storageDrivers {
		err := Init(config.StorageConfig{storageDriver, ""})
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

func TestMergeRepos(t *testing.T) {
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
			err := Driver.MergeRepos(tt.repos)
			if err != nil {
				log.Fatal("Could not merge repos:", err)
			}
			repos, err := Driver.GetRepos()
			assert.NoErr(t, err)
			assert.Equal(t, len(repos), tt.numRepos, tt.name)

			for _, r := range tt.repos {
				repo, _, err := Driver.GetRepo(*r.Name)
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
		r    *data.Repo
		want string
	}{
		{"stable repo id", &data.Repo{Name: pointerto.String("stable")}, "stable"},
		{"no id (unexpected)", &data.Repo{}, "<nil>"},
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
		r    *data.Repo
		args args
	}{
		{"stable repo id", &data.Repo{}, args{"stable"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.SetModelId(tt.args.name)
			assert.Equal(t, *tt.r.Name, tt.args.name, tt.name)
		})
	}
}

func teardownTestRepoCache() {
	if _, err := Driver.DeleteRepos(); err != nil {
		log.Fatal("Could not clear cache ", err)
	}
}
