package repos

import (
	"testing"

	"github.com/arschles/assert"
)

func TestOfficial(t *testing.T) {
	offRepo := []string{"stable", "incubator"}
	for i, repo := range official {
		assert.Equal(t, repo.Name, offRepo[i], "It contains stable")
	}
}

func TestEnabledReturnsOfficial(t *testing.T) {
	repos, err := Enabled()
	assert.NoErr(t, err)
	assert.Equal(t, repos, official, "It returns the official repos")
}
