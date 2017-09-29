package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arschles/assert"
)

func TestGetConfig(t *testing.T) {
	currentConfig = Configuration{}
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

func TestBaseDirWithEnvVar(t *testing.T) {
	path := "/path/to/monocular/home"
	os.Setenv("MONOCULAR_HOME", path)
	defer func() { os.Unsetenv("MONOCULAR_HOME") }()
	assert.Equal(t, BaseDir(), path, "BaseDir uses MONOCULAR_HOME value")
}

func TestBaseDirWithEmptyEnvVar(t *testing.T) {
	path := ""
	os.Setenv("MONOCULAR_HOME", path)
	defer func() { os.Unsetenv("MONOCULAR_HOME") }()
	assert.Equal(t, BaseDir(), path, "BaseDir uses MONOCULAR_HOME value")
}

func TestConfigFile(t *testing.T) {
	path := filepath.Join(BaseDir(), "config", "monocular.yaml")
	assert.Equal(t, configFile(), path, "Config file = BaseDir + config + monocular.yaml")
}

func TestGetOAuthConfig(t *testing.T) {
	tests := []struct {
		name            string
		clientIDSet     bool
		clientSecretSet bool
		wantErr         bool
	}{
		{"no client ID or secret", false, false, true},
		{"no client ID", false, true, true},
		{"no client secret", true, false, true},
		{"both client ID and secret", true, true, false},
	}
	expectedClientID := "clientID"
	expectedClientSecret := "clientSecret"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.clientIDSet {
				os.Setenv("MONOCULAR_AUTH_GITHUB_CLIENT_ID", expectedClientID)
				defer os.Unsetenv("MONOCULAR_AUTH_GITHUB_CLIENT_ID")
			}
			if tt.clientSecretSet {
				os.Setenv("MONOCULAR_AUTH_GITHUB_CLIENT_SECRET", expectedClientSecret)
				defer os.Unsetenv("MONOCULAR_AUTH_GITHUB_CLIENT_SECRET")
			}
			got, err := GetOAuthConfig("monocular.local")
			if tt.wantErr {
				assert.ExistsErr(t, err, tt.name)
			} else {
				assert.NotNil(t, got, "oauth config")
				assert.Equal(t, got.ClientID, expectedClientID, "client ID")
				assert.Equal(t, got.ClientSecret, expectedClientSecret, "client secret")
			}
		})
	}
}

func TestGetAuthSigningKey(t *testing.T) {
	tests := []struct {
		name    string
		set     bool
		want    string
		wantErr bool
	}{
		{"no signing key", false, "", true},
		{"signing key", true, "secret", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.set {
				os.Setenv("MONOCULAR_AUTH_SIGNING_KEY", tt.want)
				defer os.Unsetenv("MONOCULAR_AUTH_SIGNING_KEY")
			}
			got, err := GetAuthSigningKey()
			if tt.wantErr {
				assert.ExistsErr(t, err, tt.name)
			}
			assert.Equal(t, got, tt.want, "signing key")
		})
	}
}
