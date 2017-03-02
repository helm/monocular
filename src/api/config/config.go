package config

import (
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/helm/monocular/src/api/config/cors"
	"github.com/helm/monocular/src/api/config/repos"
)

// ConfigurationWithOverrides includes default Configuration values
// and environment specific ones
type configurationWithOverrides map[string]Configuration

// Configuration is the the resulting environment based Configuration
// For now it only includes Cors info
type Configuration struct {
	Cors        cors.Cors
	Repos       repos.Repos
	Initialized bool
}

// Cached version of the config
var currentConfig Configuration

// GetConfig returns the environment specific configuration
func GetConfig() (Configuration, error) {
	// Cached config
	if currentConfig.Initialized {
		return currentConfig, nil
	}

	configFilePath := configFile()

	log.WithFields(log.Fields{
		"configFile": configFilePath,
	}).Info("Configuration bootstrap init")

	_, err := os.Stat(configFilePath)
	if err == nil {
		log.Info("Configuration file found!")
	} else {
		log.Info("Configuration file not found, using defaults")
	}

	currentConfig.Cors, err = cors.Config(configFilePath)
	if err != nil {
		return currentConfig, err
	}

	currentConfig.Repos, err = repos.Enabled(configFilePath)
	if err != nil {
		return currentConfig, err
	}

	currentConfig.Initialized = true

	log.Info("Configuration bootstrap finished")
	return currentConfig, nil
}

// BaseDir returns the location of the directory
// where the configuration files are stored
func BaseDir() string {
	return filepath.Join(os.Getenv("HOME"), "monocular")
}

func configFile() string {
	return filepath.Join(BaseDir(), "config.yaml")
}
