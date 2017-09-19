package redis

import (
	"testing"

	"github.com/arschles/assert"
)

var testDriver *Driver

func getTestDriver() (*Driver, error) {
	if testDriver == nil {
		var err error
		testDriver, err = New(defaultHost)
		if err != nil {
			return nil, err
		}
	}
	return testDriver, nil
}

func TestGetRedisPool(t *testing.T) {
	driver, err := getTestDriver()
	assert.Nil(t, err, "Error initializing driver")

	pool := driver.getRedisPool()
	assert.NotNil(t, pool, "Redis Pool")
}
