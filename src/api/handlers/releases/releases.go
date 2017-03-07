package releases

import (
	"fmt"
	"net/http"
	"strings"

	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/helm/monocular/src/api/data"
	"github.com/helm/monocular/src/api/data/cache/charthelper"
	"github.com/helm/monocular/src/api/data/helpers"
	"github.com/helm/monocular/src/api/handlers"
	"github.com/helm/monocular/src/api/swagger/models"
	releasesapi "github.com/helm/monocular/src/api/swagger/restapi/operations/releases"
	hapi_release5 "k8s.io/helm/pkg/proto/hapi/release"
	rls "k8s.io/helm/pkg/proto/hapi/services"
	"k8s.io/helm/pkg/timeconv"
)

// GetReleases returns all the existing releases in your cluster
func GetReleases(helmclient data.Client, params releasesapi.GetAllReleasesParams) middleware.Responder {
	releases, err := helmclient.ListReleases(params)
	if err != nil {
		return error(err.Error())
	}

	resources := makeReleaseResources(releases)
	payload := handlers.DataResourcesBody(resources)
	return releasesapi.NewGetAllReleasesOK().WithPayload(payload)
}

// CreateRelease installs a chart version
func CreateRelease(helmclient data.Client, params releasesapi.CreateReleaseParams, c data.Charts) middleware.Responder {
	// Params validation
	format := strfmt.NewFormats()
	err := params.Data.Validate(format)
	if err != nil {
		return badRequestError(err.Error())
	}

	idSplit := strings.Split(*params.Data.ChartID, "/")
	if len(idSplit) != 2 || idSplit[0] == "" || idSplit[1] == "" {
		return badRequestError("chartId must include the repository name. i.e: stable/wordpress")
	}

	// Search chart package and get local path
	repo, chartName := idSplit[0], idSplit[1]
	chartPackage, err := c.ChartVersionFromRepo(repo, chartName, *params.Data.ChartVersion)
	if err != nil {
		return badRequestError("chart not found")
	}
	chartPath := charthelper.TarballPath(chartPackage)

	release, err := helmclient.InstallRelease(chartPath, params)
	if err != nil {
		return error(fmt.Sprintf("Can't create the release: %s", err))
	}

	resource := makeReleaseResource(release.Release)
	payload := handlers.DataResourceBody(resource)
	return releasesapi.NewCreateReleaseCreated().WithPayload(payload)
}

// DeleteRelease deletes an existing release
func DeleteRelease(helmclient data.Client, params releasesapi.DeleteReleaseParams) middleware.Responder {
	release, err := helmclient.DeleteRelease(params.ReleaseName)
	if err != nil {
		return badRequestError(fmt.Sprintf("Can't delete the release: %s", err))
	}
	resource := makeReleaseResource(release.Release)
	payload := handlers.DataResourceBody(resource)
	return releasesapi.NewDeleteReleaseOK().WithPayload(payload)
}

// error is a convenience that contains a swagger-friendly 500 given a string
func error(message string) middleware.Responder {
	return releasesapi.NewGetAllReleasesDefault(http.StatusInternalServerError).WithPayload(
		&models.Error{Code: helpers.Int64ToPtr(http.StatusInternalServerError), Message: &message},
	)
}

// error is a convenience that contains a swagger-friendly 500 given a string
func badRequestError(message string) middleware.Responder {
	return releasesapi.NewGetAllReleasesDefault(http.StatusBadRequest).WithPayload(
		&models.Error{Code: helpers.Int64ToPtr(http.StatusBadRequest), Message: &message},
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
	if release == nil {
		return &ret
	}
	ret.Type = helpers.StrToPtr("release")
	ret.ID = helpers.StrToPtr(release.Name)
	ret.Attributes = &models.Release{
		ChartName:    &release.Chart.Metadata.Name,
		ChartVersion: &release.Chart.Metadata.Version,
		Updated:      helpers.StrToPtr(timeconv.String(release.Info.LastDeployed)),
		Name:         &release.Name,
		Namespace:    &release.Namespace,
		Status:       helpers.StrToPtr(release.Info.Status.Code.String()),
	}
	return &ret
}
