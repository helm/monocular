package mocks

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/arschles/assert"
)

func TestGetYAML(t *testing.T) {
	path, err := getMocksWd()
	assert.NoErr(t, err)
	path += fmt.Sprintf("repo-%s.yaml", repoName)
	_, err = getYAML(path)
	assert.NoErr(t, err)
	path, err = getMocksWd()
	assert.NoErr(t, err)
	path += "notyaml.xml"
	_, err = getYAML(path)
	assert.ExistsErr(t, err, "confused xml for yaml, whoa")
}

func TestGetMocksWd(t *testing.T) {
	path, err := getMocksWd()
	assert.NoErr(t, err)
	assertMocksDir(t, path)
	err = os.Chdir("..") // change to the parent "api" directory
	assert.NoErr(t, err)
	path, err = getMocksWd()
	assert.NoErr(t, err)
	assertMocksDir(t, path)
	err = os.Chdir("/")
	assert.NoErr(t, err)
	_, err = getMocksWd()
	assert.ExistsErr(t, err, "can't determine mocks dir if not in 'api' dir or one of its descendant dirs")
}

func assertMocksDir(t *testing.T, path string) {
	cwdSplit := strings.Split(path, "/")
	dirname := cwdSplit[len(cwdSplit)-2]
	assert.Equal(t, dirname, "mocks", "getMocksWd a path that ends in '/mocks/'")
}
