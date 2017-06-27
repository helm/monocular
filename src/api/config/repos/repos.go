package repos

import (
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"

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
		URL:    "https://kubernetes-charts.storage.googleapis.com",
		Source: "https://github.com/kubernetes/charts/tree/master/stable",
	},
	Repo{
		Name:   "incubator",
		URL:    "https://kubernetes-charts-incubator.storage.googleapis.com",
		Source: "https://github.com/kubernetes/charts/tree/master/incubator",
	},
}

// Enabled returns the map of repositories
func Enabled(configFile string) (Repos, error) {
	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		log.Info("Loading default repositories")
		return official, nil
	}

	log.Info("Loading repositories from config file")
	repos, err := loadReposFromFile(configFile)
	if err != nil {
		return nil, err
	}

	if len(repos) == 0 {
		log.Info("No repositories found, using defaults")
		return official, nil
	}

	return repos, nil
}

func loadReposFromFile(filePath string) (Repos, error) {
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
