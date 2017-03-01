package repos

// Repos is an array of Repo
type Repos []Repo

// Repo is a map name => URL
type Repo struct {
	Name        string
	RegistryURL string
	SourceURL   string
}

var official = Repos{
	Repo{
		Name:        "stable",
		RegistryURL: "http://storage.googleapis.com/kubernetes-charts",
		SourceURL:   "https://github.com/kubernetes/charts/tree/master/stable",
	},
	Repo{
		Name:        "incubator",
		RegistryURL: "http://storage.googleapis.com/kubernetes-charts-incubator",
		SourceURL:   "https://github.com/kubernetes/charts/tree/master/incubator",
	},
}

// Enabled returns the map of repositories
func Enabled() (Repos, error) {
	// TODO, we should be able to override this from a file
	return official, nil
}
