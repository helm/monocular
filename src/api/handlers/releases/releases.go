package releases

import (
	middleware "github.com/go-openapi/runtime/middleware"
	releasesapi "github.com/helm/monocular/src/api/swagger/restapi/operations/releases"
)

// GetReleases returns all the existing releases in your cluster
func GetReleases(params releasesapi.GetAllReleasesParams) middleware.Responder {
	return releasesapi.NewGetAllReleasesOK()
}

// CreateRelease installs a chart version
func CreateRelease(params releasesapi.CreateReleaseParams) middleware.Responder {
	return releasesapi.NewCreateReleaseCreated()
}
