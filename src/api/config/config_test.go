package config

import (
	"sort"
	"testing"

	"github.com/arschles/assert"
)

func TestCurrentEnvironmentDefault(t *testing.T) {
	env := currentEnvironment()
	assert.Equal(t, env, "production", "production by default env")
}

func TestReadConfigWithOverrides(t *testing.T) {
	config := readConfigWithOverrides()
	var keys []string
	for k := range config {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var envs = []string{"default", "development"}

	for i := range keys {
		assert.Equal(t, envs[i], keys[i], "Contains the env")
	}
}

func TestGetConfig(t *testing.T) {
	config, err := GetConfig()
	assert.NoErr(t, err)
	var headers = []string{"access-control-allow-headers", "x-xsrf-token"}
	var origin = []string{"my-api-server"}
	assert.Equal(t, config.Cors.AllowedHeaders, headers, "Allowrd headers")
	assert.Equal(t, config.Cors.AllowedOrigins, origin, "Default origin")
}
