package models

import (
	"github.com/kubernetes-helm/monocular/src/api/datastore"
	"gopkg.in/mgo.v2/bson"
)

// Repo describes a chart repository
type Repo struct {
	ID     bson.ObjectId `json:"-" bson:"_id,omitempty"`
	Name   string        `json:"name" valid:"alpha,required"`
	URL    string        `json:"url" valid:"url,required"`
	Source string        `json:"source"`
}

// ReposCollection is the name of the repos collection
const ReposCollection = "repos"

// OfficialRepos are the official Kubernetes repos
var OfficialRepos = []*Repo{
	{
		Name:   "stable",
		URL:    "https://kubernetes-charts.storage.googleapis.com",
		Source: "https://github.com/kubernetes/charts/tree/master/stable",
	},
	{
		Name:   "incubator",
		URL:    "https://kubernetes-charts-incubator.storage.googleapis.com",
		Source: "https://github.com/kubernetes/charts/tree/master/incubator",
	},
}

// ListRepos returns a list of all repos
func ListRepos(db datastore.Database) ([]*Repo, error) {
	c := db.C(ReposCollection)
	var repos []*Repo
	err := c.Find(nil).All(&repos)
	return repos, err
}

// GetRepo by name
func GetRepo(db datastore.Database, name string) (*Repo, error) {
	c := db.C(ReposCollection)
	var repo Repo
	err := c.Find(bson.M{"name": name}).One(&repo)
	if err != nil {
		return nil, err
	}
	return &repo, nil
}

// CreateRepos takes an array of Repos and saves them to the database
// Repos with the same name are updated
func CreateRepos(db datastore.Database, repos []*Repo) error {
	for _, r := range repos {
		err := CreateRepo(db, r)
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateRepo takes a Repo object and saves it to the database
// Repos with the same name are updated
func CreateRepo(db datastore.Database, repo *Repo) error {
	c := db.C(ReposCollection)
	_, err := c.Upsert(bson.M{"name": repo.Name}, repo)
	return err
}

// DeleteRepo by name
func DeleteRepo(db datastore.Database, name string) error {
	c := db.C(ReposCollection)
	return c.Remove(bson.M{"name": name})
}

// func (r *Repo) IndexURL() *url.URL {
// 	u, _ := url.Parse(r.URL)
// 	u.Path = path.Join(u.Path, "index.yaml")
// 	return u
// }
