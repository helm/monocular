package mocks

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/helm/monocular/src/api/data/helpers"
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

// getMocksWd returns the full local pathname of the monocular API mocks directory, including a trailing "/"
// The purpose is for various sub-packages to load filesystem-derived mocks from their cwd
func getMocksWd() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	cwdSplit := strings.Split(cwd, "/")
	// find out where we are in relation to the "api" directory
	// it's an O(n) operation, but we should never have, like, a lot of directories (crosses fingers)
	for i, dirname := range cwdSplit {
		if dirname == "api" {
			cwdSplit = cwdSplit[:i+1] // strip directories after "api"
			return strings.Join(cwdSplit, "/") + "/mocks/", nil
		}
	}
	return cwd, fmt.Errorf("can't get mocks dir unless cwd is a child of 'api' dir")
}
