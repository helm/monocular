package config

import (
	"path/filepath"
	"testing"

	"github.com/arschles/assert"
)

func TestNewRedisPool(t *testing.T) {
	currentConfig = Configuration{}
	NewRedisPool()
	assert.NotNil(t, Pool, "Redis Pool")
	err := Pool.Close()
	assert.NoErr(t, err)
}

func Test_getRedisConf(t *testing.T) {
	tests := []struct {
		name           string
		configFileName string
		expectedHost   string
	}{
		{"No Redis config", "noredis_config.yaml", defaultHost},
		{"Blank Redis config", "emptyredis_config.yaml", defaultHost},
		{"Custom Redis config", "config.yaml", "myredis:1234"},
	}
	configFileOrig := configFile
	defer func() { configFile = configFileOrig }()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentConfig = Configuration{}
			configFile = func() string {
				return filepath.Join("./testdata", tt.configFileName)
			}
			conf := getRedisConf()
			assert.Equal(t, conf.Host, tt.expectedHost, tt.name)
		})
	}
}
