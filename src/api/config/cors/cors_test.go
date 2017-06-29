package cors

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/arschles/assert"
)

var configFileOk = filepath.Join("..", "testdata", "config.yaml")
var configFileNotOk = filepath.Join("..", "testdata", "bogus_config.yaml")
var configFileNoCors = filepath.Join("..", "testdata", "nocors_config.yaml")
var defaultExpectedCors = Cors{
	AllowedOrigins: []string{""},
	AllowedHeaders: []string{"content-type", "x-xsrf-token"},
}

func TestConfigFileDoesNotExist(t *testing.T) {
	config, err := Config("no-file")
	assert.NoErr(t, err)
	assert.Equal(t, config.AllowedHeaders, defaultExpectedCors.AllowedHeaders, "Allowed headers")
	assert.Equal(t, config.AllowedOrigins, defaultExpectedCors.AllowedOrigins, "Default origin")
}

// In development environment, CORS has a permissive configuration
func TestConfigFileDoesNotExistDevelopment(t *testing.T) {
	os.Setenv("ENVIRONMENT", "development")
	defer func() { os.Unsetenv("ENVIRONMENT") }()
	var origin = []string{"*"}
	config, err := Config("no-file")
	assert.NoErr(t, err)
	assert.Equal(t, len(config.AllowedHeaders), 0, "Allowed headers")
	assert.Equal(t, config.AllowedOrigins, origin, "Default origin")
}

func TestConfigFileWithoutCors(t *testing.T) {
	cors, err := Config(configFileNoCors)
	assert.NoErr(t, err)
	assert.Equal(t, cors, defaultExpectedCors, "It returns the default CORS")
}

func TestConfigFromFile(t *testing.T) {
	expected := Cors{
		AllowedOrigins: []string{"http://mymonocular"},
		AllowedHeaders: []string{"access-control-allow-headers", "x-xsrf-token"},
	}
	cors, err := Config(configFileOk)
	assert.NoErr(t, err)
	assert.Equal(t, cors, expected, "It uses the cors from the config file")
}

// Return err
func TestConfigFromFileInvalid(t *testing.T) {
	_, err := Config(configFileNotOk)
	assert.ExistsErr(t, err, "File exist but it is not valid")
}

func TestLoadCorsFromFileDoesNotExist(t *testing.T) {
	cors, err := loadCorsFromFile("does not exist")
	assert.ExistsErr(t, err, "Can not load the file")
	assert.Equal(t, cors, Cors{}, "Returns no cors")
}
