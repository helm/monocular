package repos

import (
	"fmt"
	"os"
)

// Repos is an array of Repo
type Repos []Repo

type reposYAML struct {
	Repositories Repo
}

// Repo is a map name => URL
type Repo struct {
	Name   string
	URL    string `yaml:"registry"`
	Source string `yaml:"source"`
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
	repos := official
	fmt.Println(configFile)
	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		return repos, nil
	}
	//repos, err = loadReposFromConfigFile(configFile)
	//if err != nil {
	//	return repos, err
	//}

	return repos, nil
}

//func loadReposFromConfigFile(filePath string) (Repos, error) {
//	var repos Repos
//	var yamlStruct reposYAML
//	bytes, err := ioutil.ReadFile(filePath)
//	if err != nil {
//		return nil, err
//	}
//	if err := yaml.Unmarshal(bytes, &yamlStruct); err != nil {
//		return nil, err
//	}
//	fmt.Printf("WAPS %+v", yamlStruct.Repositories)
//	return repos, nil
//}
