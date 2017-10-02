package repos

import (
	"path/filepath"
	"testing"

	"github.com/arschles/assert"
	"github.com/kubernetes-helm/monocular/src/api/models"
)

var configFileOk = filepath.Join("..", "testdata", "config.yaml")
var configFileNotOk = filepath.Join("..", "testdata", "bogus_config.yaml")
var configFileNoRepos = filepath.Join("..", "testdata", "norepos_config.yaml")

func TestOfficial(t *testing.T) {
	offRepo := []models.Repo{
		{
			Name: "stable",
		},
		{
			Name: "incubator",
		},
	}
	for i, repo := range official {
		assert.Equal(t, repo.Name, offRepo[i].Name, "It contains only official repos")
	}
}

func TestEnabledFileDoesnotExist(t *testing.T) {
	repos, err := Enabled("no-file")
	assert.NoErr(t, err)
	assert.Equal(t, len(repos), len(official), "same length")
	for i := 0; i < len(repos); i++ {
		assert.Equal(t, repos[i], official[i], "correct repos")
	}
}

func TestEnabledFileWithoutRepos(t *testing.T) {
	repos, err := Enabled(configFileNoRepos)
	assert.NoErr(t, err)
	assert.Equal(t, len(repos), len(official), "same length")
	for i := 0; i < len(repos); i++ {
		assert.Equal(t, repos[i], official[i], "correct repos")
	}
}

// Use the repositories in the file
func TestEnabledReposInFile(t *testing.T) {
	repos, err := Enabled(configFileOk)
	assert.NoErr(t, err)
	offRepo := []models.Repo{
		{
			Name:   "repoName",
			URL:    "http://myrepobucket",
			Source: "http://github.com/my-repo",
		},
		{
			Name: "repoName2",
			URL:  "http://myrepobucket2",
		},
	}

	assert.Equal(t, len(repos), 2, "Only has repos from the YAML file")

	for i, repo := range repos {
		assert.Equal(t, repo.Name, offRepo[i].Name, "Same repo name")
		assert.Equal(t, repo.URL, offRepo[i].URL, "Same repo URL")
		assert.Equal(t, repo.Source, offRepo[i].Source, "Same repo Source")
	}
}

// Return err
func TestEnabledWrongFile(t *testing.T) {
	_, err := Enabled(configFileNotOk)
	assert.ExistsErr(t, err, "File exist but it is not valid")
}

func TestLoadReposFromFile(t *testing.T) {
	repos, err := loadReposFromFile("does not exist")
	assert.ExistsErr(t, err, "Can not load the file")
	assert.Equal(t, len(repos), 0, "Returns no repos")
}
