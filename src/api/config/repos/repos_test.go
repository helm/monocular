package repos

import (
	"path/filepath"
	"testing"

	"github.com/arschles/assert"
)

var configFileOk = filepath.Join("..", "testdata", "config.yaml")
var configFileNotOk = filepath.Join("..", "testdata", "bogus_config.yaml")

func TestOfficial(t *testing.T) {
	offRepo := []Repo{
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
	assert.Equal(t, repos, official, "It returns the official repos")
}

// Use the repositories in the file
func TestEnabledReposInFile(t *testing.T) {
	_, err := Enabled(configFileOk)
	assert.NoErr(t, err)
}

// Return err
func TestEnabledWrongFile(t *testing.T) {

}
