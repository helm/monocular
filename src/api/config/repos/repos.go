package repos

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Repos is an array of Repo
type Repos []Repo

type reposYAML struct {
	Repos Repos
}

// Repo is a map name => URL
type Repo struct {
	Name   string
	URL    string
	Source string
}

var official = Repos{
	Repo{
		Name:   "stable",
		URL:    "http://storage.googleapis.com/kubernetes-charts",
		Source: "https://github.com/kubernetes/charts/tree/master/stable",
	},
	Repo{
		Name:   "incubator",
		URL:    "http://storage.googleapis.com/kubernetes-charts-incubator",
		Source: "https://github.com/kubernetes/charts/tree/master/incubator",
	},
}

// Enabled returns the map of repositories
func Enabled(configFile string) (Repos, error) {
	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		return official, nil
	}
	repos, err := loadReposFromFile(configFile)
	if err != nil {
		return nil, err
	}

	return repos, nil
}

func loadReposFromFile(filePath string) (Repos, error) {
	fmt.Printf("Loading repos from file %s\n", filePath)
	var yamlStruct reposYAML
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(bytes, &yamlStruct); err != nil {
		return nil, err
	}
	return yamlStruct.Repos, nil
}
