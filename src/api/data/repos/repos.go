package repos

// Repos is an array of Repo
type Repos []Repo

// Repo is a map name => URL
type Repo struct {
	Name string
	URL  string
}

var official = Repos{
	Repo{
		Name: "stable",
		URL:  "http://storage.googleapis.com/kubernetes-charts",
	},
	Repo{
		Name: "incubator",
		URL:  "http://storage.googleapis.com/kubernetes-charts-incubator",
	},
}

// Enabled returns the map of repositories
func Enabled() (Repos, error) {
	// TODO, we should be able to override this from a file
	return official, nil
}
