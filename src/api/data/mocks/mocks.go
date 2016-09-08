package mocks

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/helm/monocular/src/api/pkg/swagger/models"
	"gopkg.in/yaml.v2"
)

// getYAML gets a yaml file from the local filesystem
func getYAML(filepath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Printf("Error reading YAML file %s: %#v\n", filepath, err)
		return nil, err
	}
	if !isYAML(data) {
		return nil, fmt.Errorf("data is not valid YAML")
	}
	return data, nil
}

// isYAML checks for valid YAML
func isYAML(b []byte) bool {
	var yml map[string]interface{}
	return yaml.Unmarshal(b, &yml) == nil
}

// ParseYAMLChartVersion converts a YAML representation of a versioned chart
// to a ChartVersion type
func ParseYAMLChartVersion(rawYAML []byte) (models.ChartVersion, error) {
	var chart models.ChartVersion
	if err := yaml.Unmarshal(rawYAML, &chart); err != nil {
		return models.ChartVersion{}, err
	}
	return chart, nil
}

func getMocksWd() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	cwdSplit := strings.Split(cwd, "/")
	if (cwdSplit[len(cwdSplit)-1]) != "api" {
		cwdSplit = cwdSplit[:len(cwdSplit)-1] // strip last directory
		return strings.Join(cwdSplit, "/") + "/data/mocks/"
	}
	return cwd + "/data/mocks/"
}
