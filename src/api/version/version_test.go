package version

import (
	"testing"

	"github.com/arschles/assert"
)

func TestGetUserAgent(t *testing.T) {
	oldVersion := Version
	Version = "0.5.4"
	defer func() { Version = oldVersion }()
	assert.Equal(t, GetUserAgent(), "monocular/0.5.4", "user agent")
}
