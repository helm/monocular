package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"

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
	Cors                 cors.Cors
	Repos                repos.Repos
	ReleasesEnabled      bool  `yaml:"releasesEnabled"`
	TillerPortForward    bool  `yaml:"tillerPortForward"`
	CacheRefreshInterval int64 `yaml:"cacheRefreshInterval"`
	Initialized          bool
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
		// Load custom configuration
		loadConfigFromFile(configFilePath, &currentConfig)
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

var configFile = func() string {
	return filepath.Join(BaseDir(), "config", "monocular.yaml")
}

func loadConfigFromFile(filePath string, configStruct *Configuration) error {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(bytes, configStruct); err != nil {
		return err
	}
	return nil
}
