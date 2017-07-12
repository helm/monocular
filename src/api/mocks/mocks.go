package mocks

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/kubernetes-helm/monocular/src/api/data/helpers"
)

// getYAML gets a yaml file from the local filesystem
func getYAML(filepath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Printf("Error reading YAML file %s: %#v\n", filepath, err)
		return nil, err
	}
	if !helpers.IsYAML(data) {
		return nil, fmt.Errorf("data is not valid YAML")
	}
	return data, nil
}

// getTestDataWd returns the full local pathname of the monocular API mocks/testdata directory, including a trailing "/"
// The purpose is for various sub-packages to load filesystem-derived mocks from their cwd
// Requires a cwd of either the monocular git repo root, or somewhere inside 'src/api', inclusive
var getTestDataWd = func() (string, error) {
	const apiDir = "/src/api"
	const testdataDir = "/mocks/testdata/"
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	cwdSplit := strings.Split(cwd, "/")
	// are we in the repo root dir?
	if cwdSplit[len(cwdSplit)-1] == "monocular" {
		return cwd + apiDir + testdataDir, nil
	}
	// find out where we are in relation to the "api" directory
	// it's an O(n) operation, but we should never have, like, a lot of directories (crosses fingers)
	for i, dirname := range cwdSplit {
		if dirname == "api" {
			cwdSplit = cwdSplit[:i+1] // strip directories after "api"
			return strings.Join(cwdSplit, "/") + testdataDir, nil
		}
	}
	return cwd, fmt.Errorf("couldn't locate ourselves in relation to the 'src/api' directory of the monocular repo")
}
