package redis

import (
	"testing"

	"github.com/arschles/assert"
)

var testDriver *Driver

func getTestDriver() *Driver {
	if testDriver == nil {
		testDriver = New(defaultHost)
	}
	return testDriver
}

func TestGetRedisPool(t *testing.T) {
	driver := getTestDriver()
	pool := driver.getRedisPool()
	assert.NotNil(t, pool, "Redis Pool")
}
