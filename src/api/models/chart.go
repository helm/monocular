package models

import (
	"github.com/kubernetes-helm/monocular/src/api/datastore"
	"gopkg.in/mgo.v2/bson"
)

// Chart describes a chart
type Chart struct {
	ID         bson.ObjectId   `json:"-" bson:"_id,omitempty"`
	Name       string          `json:"name"`
	Repo       string          `json:"repo"`
	stargazers []bson.ObjectId `bson:"stargazers"`
}

// ChartsCollection is the name of the charts collection
const ChartsCollection = "charts"

// CreateChart takes a Chart object and saves it to the database
// Charts with the same name and repo are updated
func CreateChart(db datastore.Database, chart *Chart) error {
	c := db.C(ChartsCollection)
	_, err := c.Upsert(bson.M{"name": chart.Name, "repo": chart.Repo}, chart)
	return err
}

// GetChartByName returns the chart matching the repo and chart name
func GetChartByName(db datastore.Database, repo string, chartName string) (*Chart, error) {
	c := db.C(ChartsCollection)
	var chart Chart
	if err := c.Find(bson.M{"name": chartName, "repo": repo}).One(&chart); err != nil {
		return nil, err
	}
	return &chart, nil
}
