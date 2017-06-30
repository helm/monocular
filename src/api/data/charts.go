package data

import (
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
	"github.com/kubernetes-helm/monocular/src/api/swagger/restapi/operations/charts"
)

// Charts is an interface for managing chart data sourced from a repository index
type Charts interface {
	// ChartFromRepo retrieves the latest version of a particular chart from a repo
	ChartFromRepo(repo, name string) (*models.ChartPackage, error)
	// ChartVersionFromRepo retrieves a specific chart version from a repo
	ChartVersionFromRepo(repo, name, version string) (*models.ChartPackage, error)
	// ChartVersionsFromRepo retrieves all chart versions from a repo
	ChartVersionsFromRepo(repo, name string) ([]*models.ChartPackage, error)
	// AllFromRepo retrieves all charts from a repo
	AllFromRepo(repo string) ([]*models.ChartPackage, error)
	// All retrieves all charts from all repos
	All() ([]*models.ChartPackage, error)
	// Search operates against all charts/repos
	Search(params charts.SearchChartsParams) ([]*models.ChartPackage, error)
	// Refresh freshens charts data
	Refresh() error
}
