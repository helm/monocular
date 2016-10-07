package data

import "github.com/helm/monocular/src/api/swagger/models"

// Charts is an interface for managing chart data
type Charts interface {
	// will have a ChartFromRepo method to retrieve a particular chart from a repo
	ChartFromRepo(repo, name string) (*models.ChartVersion, error)
	// will have a AllFromRepo method to retrieve all charts from a repo
	AllFromRepo(repo string) ([]*models.ChartVersion, error)
	// will have a All method to retrieve all charts from all repos
	All() ([]*models.Resource, error)
}
