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

func getMocksWd() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	cwdSplit := strings.Split(cwd, "/")
	if (cwdSplit[len(cwdSplit)-1]) != "api" {
		cwdSplit = cwdSplit[:len(cwdSplit)-1] // strip last directory
		return strings.Join(cwdSplit, "/") + "/mocks/"
	}
	return cwd + "/mocks/"
}
