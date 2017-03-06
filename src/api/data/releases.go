package data

import (
	releasesapi "github.com/helm/monocular/src/api/swagger/restapi/operations/releases"
	rls "k8s.io/helm/pkg/proto/hapi/services"
)

// Releases is an interface for managing Helm Chart releases
type Releases interface {
	// ListReleases retrieves the list of Helm releases deployed in your cluster
	ListReleases() (*rls.ListReleasesResponse, error)
	// InstallRelease creates a new released based on an existing chartPackage
	InstallRelease(chartPath string, params releasesapi.CreateReleaseParams) (*rls.InstallReleaseResponse, error)
}
