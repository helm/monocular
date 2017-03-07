package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arschles/assert"
)

func TestGetConfig(t *testing.T) {
	configFileOrig := configFile
	defer func() { configFile = configFileOrig }()
	configFile = func() string {
		return filepath.Join("./testdata", "doesnotexist.yaml")
	}
	config, err := GetConfig()
	assert.NoErr(t, err)
	if len(config.Cors.AllowedHeaders) == 0 {
		t.Error("AllowedHeaders not present")
	}
	if len(config.Cors.AllowedOrigins) == 0 {
		t.Error("AllowedOrigins not present")
	}
	if len(config.Repos) == 0 {
		t.Error("Repositories not present")
	}

	assert.False(t, config.ReleasesEnabled, "Releases disabled")
}

func TestGetConfigFromFile(t *testing.T) {
	currentConfig = Configuration{}
	configFileOrig := configFile
	defer func() { configFile = configFileOrig }()
	configFile = func() string {
		return filepath.Join("./testdata", "config.yaml")
	}
	config, err := GetConfig()
	assert.NoErr(t, err)
	if len(config.Cors.AllowedHeaders) == 0 {
		t.Error("AllowedHeaders not present")
	}
	if len(config.Cors.AllowedOrigins) == 0 {
		t.Error("AllowedOrigins not present")
	}
	if len(config.Repos) == 0 {
		t.Error("Repositories not present")
	}
	assert.Equal(t, config.Repos[0].Name, "repoName", "First repo")
	assert.Equal(t, config.Repos[1].Name, "repoName2", "Second repo")
	assert.Equal(t, config.Repos[0].URL, "http://myrepobucket", "Repo URL")
	assert.Equal(t, config.Repos[1].URL, "http://myrepobucket2", "Repo URL")
	assert.Equal(t, config.Repos[0].Source, "http://github.com/my-repo", "Repo Source")
	assert.Equal(t, config.Repos[1].Source, "", "Repo Source")

	assert.True(t, config.ReleasesEnabled, "Releases disabled")
}

func TestLoadConfigFromFile(t *testing.T) {
	config := Configuration{}
	err := loadConfigFromFile("./testdata/config.yaml", &config)
	assert.NoErr(t, err)
	err = loadConfigFromFile("./testdata/does-not-exist.yaml", &config)
	assert.ExistsErr(t, err, "File does not exist")
	err = loadConfigFromFile("./config.go", &config)
	assert.ExistsErr(t, err, "File not valid")
}

func TestBaseDir(t *testing.T) {
	path := filepath.Join(os.Getenv("HOME"), "monocular")
	assert.Equal(t, BaseDir(), path, "BaseDir uses home + monocular")
}

func TestConfigFile(t *testing.T) {
	path := filepath.Join(BaseDir(), "config", "monocular.yaml")
	assert.Equal(t, configFile(), path, "Config file = BaseDir + config + monocular.yaml")
}
