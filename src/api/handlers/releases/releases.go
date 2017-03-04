package releases

import (
	"fmt"
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/helm/monocular/src/api/swagger/models"
	releasesapi "github.com/helm/monocular/src/api/swagger/restapi/operations/releases"

	helmclient "github.com/helm/monocular/src/api/data/helm/client"
	helmreleases "github.com/helm/monocular/src/api/data/helm/releases"
	"github.com/helm/monocular/src/api/data/helpers"
)

// GetReleases returns all the existing releases in your cluster
func GetReleases(params releasesapi.GetAllReleasesParams) middleware.Responder {
	client, err := helmclient.CreateTillerClient()
	if err != nil {
		return error("Error creating the Helm client")
	}
	releases, err := helmreleases.ListReleases(client)
	if err != nil {
		return error("Error retrieving the list of releases")
	}
	fmt.Printf("%+v", releases)
	return releasesapi.NewGetAllReleasesOK()
}

// CreateRelease installs a chart version
func CreateRelease(params releasesapi.CreateReleaseParams) middleware.Responder {
	return releasesapi.NewCreateReleaseCreated()
}

// error is a convenience that contains a swagger-friendly 500 given a resource string
func error(resource string) middleware.Responder {
	message := fmt.Sprintf("500 %s not found", resource)
	return releasesapi.NewGetAllReleasesDefault(http.StatusExpectationFailed).WithPayload(
		&models.Error{Code: helpers.Int64ToPtr(http.StatusExpectationFailed), Message: &message},
	)
}
