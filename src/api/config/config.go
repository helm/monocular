package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	oauth2Github "golang.org/x/oauth2/github"
	yaml "gopkg.in/yaml.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/kubernetes-helm/monocular/src/api/config/cors"
	"github.com/kubernetes-helm/monocular/src/api/config/repos"
	"github.com/kubernetes-helm/monocular/src/api/datastore"
)

// ConfigurationWithOverrides includes default Configuration values
// and environment specific ones
type configurationWithOverrides map[string]Configuration

type oauthConfig struct {
	ClientID     string `yaml:"clientID"`
	ClientSecret string `yaml:"clientSecret"`
}

// Configuration is the the resulting environment based Configuration
// For now it only includes Cors info
type Configuration struct {
	Cors                 cors.Cors
	Repos                repos.Repos
	ReleasesEnabled      bool             `yaml:"releasesEnabled"`
	TillerPortForward    bool             `yaml:"tillerPortForward"`
	CacheRefreshInterval int64            `yaml:"cacheRefreshInterval"`
	Mongo                datastore.Config `yaml:"mongodb"`
	OAuthConfig          oauthConfig      `yaml:"oauthConfig"`
	SigningKey           string           `yaml:"signingKey"`
	TillerNamespace      string           `yaml:"tillerNamespace"`
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

	if currentConfig.TillerNamespace == "" {
		currentConfig.TillerNamespace = "kube-system"
	}

	currentConfig.Initialized = true

	log.Info("Configuration bootstrap finished")
	return currentConfig, nil
}

// GetOAuthConfig returns the OAuth configuration for the GitHub provider
func GetOAuthConfig(host string) (*oauth2.Config, error) {
	clientID, ok := os.LookupEnv("MONOCULAR_AUTH_GITHUB_CLIENT_ID")
	if !ok {
		return nil, errors.New("no client ID for GitHub provider, ensure MONOCULAR_AUTH_GITHUB_CLIENT_ID is set")
	}
	clientSecret, ok := os.LookupEnv("MONOCULAR_AUTH_GITHUB_CLIENT_SECRET")
	if !ok {
		return nil, errors.New("no client secret for GitHub provider, ensure MONOCULAR_AUTH_GITHUB_CLIENT_SECRET is set")
	}
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     oauth2Github.Endpoint,
		RedirectURL:  "http://" + host + "/api/auth/github/callback",
		Scopes:       []string{"user:email"},
	}, nil
}

// GetAuthSigningKey returns the secret key used for signing JWTs and secure sessions
func GetAuthSigningKey() (string, error) {
	signingKey, ok := os.LookupEnv("MONOCULAR_AUTH_SIGNING_KEY")
	if !ok {
		return "", errors.New("no signing key, ensure MONOCULAR_AUTH_SIGNING_KEY is set")
	}
	return signingKey, nil
}

// BaseDir returns the location of the directory
// where the configuration files are stored
func BaseDir() string {
	if basedir, ok := os.LookupEnv("MONOCULAR_HOME"); ok {
		return basedir
	}

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
