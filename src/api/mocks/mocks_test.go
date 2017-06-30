package mocks

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/arschles/assert"
	"github.com/kubernetes-helm/monocular/src/api/testutil"
)

func TestGetYAML(t *testing.T) {
	path, err := getTestDataWd()
	assert.NoErr(t, err)
	path += fmt.Sprintf("repo-%s.yaml", testutil.RepoName)
	_, err = getYAML(path)
	assert.NoErr(t, err)
	path, err = getTestDataWd()
	assert.NoErr(t, err)
	path += "notyaml.xml"
	_, err = getYAML(path)
	assert.ExistsErr(t, err, "confused xml for yaml, whoa")
}

func TestGetMocksWd(t *testing.T) {
	path, err := getTestDataWd()
	assert.NoErr(t, err)
	assertMocksDir(t, path)
	err = os.Chdir("..") // change to the parent "api" directory
	assert.NoErr(t, err)
	path, err = getTestDataWd()
	assert.NoErr(t, err)
	assertMocksDir(t, path)
	err = os.Chdir("../..") // change to the git repo root directory
	assert.NoErr(t, err)
	path, err = getTestDataWd()
	assert.NoErr(t, err)
	assertMocksDir(t, path)
	err = os.Chdir("/")
	assert.NoErr(t, err)
	_, err = getTestDataWd()
	assert.ExistsErr(t, err, "couldn't locate ourselves in relation to the 'src/api' directory of the monocular repo")
}

func assertMocksDir(t *testing.T, path string) {
	cwdSplit := strings.Split(path, "/")
	dirname := cwdSplit[len(cwdSplit)-2]
	assert.Equal(t, dirname, "testdata", "getTestDataWd a path that ends in '/testdata/'")
}
