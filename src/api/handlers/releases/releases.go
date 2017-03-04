package releases

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
	helmclient "github.com/helm/monocular/src/api/data/helm/client"
	helmreleases "github.com/helm/monocular/src/api/data/helm/releases"
	"github.com/helm/monocular/src/api/data/helpers"
	"github.com/helm/monocular/src/api/handlers"
	"github.com/helm/monocular/src/api/swagger/models"
	releasesapi "github.com/helm/monocular/src/api/swagger/restapi/operations/releases"
	hapi_release5 "k8s.io/helm/pkg/proto/hapi/release"
	rls "k8s.io/helm/pkg/proto/hapi/services"
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

	resources := makeReleaseResources(releases)
	payload := handlers.DataResourcesBody(resources)
	return releasesapi.NewGetAllReleasesOK().WithPayload(payload)
}

// CreateRelease installs a chart version
func CreateRelease(params releasesapi.CreateReleaseParams) middleware.Responder {
	return releasesapi.NewCreateReleaseCreated()
}

// error is a convenience that contains a swagger-friendly 500 given a string
func error(message string) middleware.Responder {
	return releasesapi.NewGetAllReleasesDefault(http.StatusInternalServerError).WithPayload(
		&models.Error{Code: helpers.Int64ToPtr(http.StatusInternalServerError), Message: &message},
	)
}

func makeReleaseResources(releases *rls.ListReleasesResponse) []*models.Resource {
	var resources []*models.Resource
	for _, release := range releases.Releases {
		resource := makeReleaseResource(release)
		resources = append(resources, resource)
	}
	return resources
}

func makeReleaseResource(release *hapi_release5.Release) *models.Resource {
	var ret models.Resource
	ret.Type = helpers.StrToPtr("release")
	ret.ID = helpers.StrToPtr(release.Name)
	ret.Attributes = &models.Release{
		Chart:     &release.Chart.Metadata.Name,
		Name:      &release.Name,
		Namespace: &release.Namespace,
		Status:    helpers.StrToPtr(release.Info.Status.Code.String()),
	}
	return &ret
}
