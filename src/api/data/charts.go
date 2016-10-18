package data

import "github.com/helm/monocular/src/api/swagger/models"

// Charts is an interface for managing chart data
type Charts interface {
	// ChartFromRepo retrieves the latest version of a particular chart from a repo
	ChartFromRepo(repo, name string) (*models.ChartVersion, error)
	// ChartVersionFromRepo retrieves a specific chart version from a repo
	ChartVersionFromRepo(repo, name, version string) (*models.ChartVersion, error)
	// ChartVersionsFromRepo retrieves all chart versions from a repo
	ChartVersionsFromRepo(repo, name string) ([]*models.ChartVersion, error)
	// AllFromRepo retrieves all charts from a repo
	AllFromRepo(repo string) ([]*models.ChartVersion, error)
	// All retrieves all charts from all repos
	All() ([]*models.Resource, error)
	// Refresh freshens charts data
	Refresh() error
}
