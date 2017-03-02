package config

import (
	"os"
	"path/filepath"

	"github.com/helm/monocular/src/api/config/cors"
	"github.com/helm/monocular/src/api/config/repos"
)

// ConfigurationWithOverrides includes default Configuration values
// and environment specific ones
type configurationWithOverrides map[string]Configuration

// Configuration is the the resulting environment based Configuration
// For now it only includes Cors info
type Configuration struct {
	Cors  cors.Cors
	Repos repos.Repos
}

// GetConfig returns the environment specific configuration
func GetConfig() (Configuration, error) {
	res := Configuration{}
	var err error
	res.Cors, err = cors.Config(configFile())
	if err != nil {
		return res, err
	}
	res.Repos, err = repos.Enabled(configFile())
	if err != nil {
		return res, err
	}

	return res, nil
}

// BaseDir returns the location of the directory
// where the configuration files are stored
func BaseDir() string {
	return filepath.Join(os.Getenv("HOME"), "monocular")
}

func configFile() string {
	return filepath.Join(BaseDir(), "config.yaml")
}
