package repos

import (
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/kubernetes-helm/monocular/src/api/models"

	yaml "gopkg.in/yaml.v2"
)

// Repos is an array of models.Repo
type Repos []*models.Repo

type reposYAML struct {
	Repos Repos
}

var official = models.OfficialRepos

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
