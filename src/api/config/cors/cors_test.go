package cors

import (
	"testing"

	"github.com/arschles/assert"
)

func TestConfig(t *testing.T) {
	var headers = []string{"access-control-allow-headers", "x-xsrf-token"}
	var origin = []string{"my-api-server"}
	config, err := Config("no-file")
	assert.NoErr(t, err)
	assert.Equal(t, config.AllowedHeaders, headers, "Allowed headers")
	assert.Equal(t, config.AllowedOrigins, origin, "Default origin")
}

// In development environment, CORS has a permissive configuration
func TestConfigDevelopment(t *testing.T) {
	origCurrentEnv := currentEnv
	currentEnv = func() string {
		return "development"
	}
	defer func() { currentEnv = origCurrentEnv }()
	var origin = []string{"*"}
	config, err := Config("no-file")
	assert.NoErr(t, err)
	assert.Equal(t, len(config.AllowedHeaders), 0, "Allowed headers")
	assert.Equal(t, config.AllowedOrigins, origin, "Default origin")
}
