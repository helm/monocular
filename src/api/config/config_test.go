package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arschles/assert"
)

func TestGetConfig(t *testing.T) {
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
}

func TestBaseDir(t *testing.T) {
	path := filepath.Join(os.Getenv("HOME"), "monocular")
	assert.Equal(t, BaseDir(), path, "BaseDir uses home + monocular")
}

func TestConfigFile(t *testing.T) {
	path := filepath.Join(BaseDir(), "config.yaml")
	assert.Equal(t, configFile(), path, "Config file = BaseDir + config.yaml")
}
