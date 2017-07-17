package releases

import (
	"fmt"
	"net/http"
	"strings"

	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/kubernetes-helm/monocular/src/api/data"
	"github.com/kubernetes-helm/monocular/src/api/data/cache/charthelper"
	"github.com/kubernetes-helm/monocular/src/api/data/helpers"
	"github.com/kubernetes-helm/monocular/src/api/handlers"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
	releasesapi "github.com/kubernetes-helm/monocular/src/api/swagger/restapi/operations/releases"
	hapi_release5 "k8s.io/helm/pkg/proto/hapi/release"
	rls "k8s.io/helm/pkg/proto/hapi/services"
	"k8s.io/helm/pkg/timeconv"
)

// GetReleases returns all the existing releases in your cluster
func GetReleases(helmclient data.Client, params releasesapi.GetAllReleasesParams, releasesEnabled bool) middleware.Responder {
	if !releasesEnabled {
		return errorResponse("Feature not enabled", http.StatusForbidden)
	}

	releases, err := helmclient.ListReleases(params)
	if err != nil {
		return errorResponse(err.Error(), http.StatusInternalServerError)
	}

	resources := makeReleaseResources(releases)
	payload := handlers.DataResourcesBody(resources)
	return releasesapi.NewGetAllReleasesOK().WithPayload(payload)
}

// GetRelease returns the extended version of a release
func GetRelease(helmclient data.Client, params releasesapi.GetReleaseParams, releasesEnabled bool) middleware.Responder {
	if !releasesEnabled {
		return errorResponse("Feature not enabled", http.StatusForbidden)
	}

	release, err := helmclient.GetRelease(params.ReleaseName)
	if err != nil {
		return errorResponse(err.Error(), http.StatusInternalServerError)
	}

	resource := makeReleaseExtendedResource(release.Release)
	payload := handlers.DataResourceBody(resource)
	return releasesapi.NewGetReleaseOK().WithPayload(payload)
}

// CreateRelease installs a chart version
func CreateRelease(helmclient data.Client, params releasesapi.CreateReleaseParams, c data.Charts, releasesEnabled bool) middleware.Responder {
	if !releasesEnabled {
		return errorResponse("Feature not enabled", http.StatusForbidden)
	}

	// Params validation
	format := strfmt.NewFormats()
	err := params.Data.Validate(format)
	if err != nil {
		return errorResponse(err.Error(), http.StatusBadRequest)
	}

	idSplit := strings.Split(*params.Data.ChartID, "/")
	if len(idSplit) != 2 || idSplit[0] == "" || idSplit[1] == "" {
		return errorResponse("chartId must include the repository name. i.e: stable/wordpress", http.StatusBadRequest)
	}

	// Search chart package and get local path
	repo, chartName := idSplit[0], idSplit[1]
	chartPackage, err := c.ChartVersionFromRepo(repo, chartName, *params.Data.ChartVersion)
	if err != nil {
		return errorResponse("chart not found", http.StatusBadRequest)
	}
	chartPath := charthelper.TarballPath(chartPackage)

	release, err := helmclient.InstallRelease(chartPath, params)
	if err != nil {
		return errorResponse(fmt.Sprintf("Can't create the release: %s", err), http.StatusInternalServerError)
	}

	resource := makeReleaseResource(release.Release)
	payload := handlers.DataResourceBody(resource)
	return releasesapi.NewCreateReleaseCreated().WithPayload(payload)
}

// DeleteRelease deletes an existing release
func DeleteRelease(helmclient data.Client, params releasesapi.DeleteReleaseParams, releasesEnabled bool) middleware.Responder {
	if !releasesEnabled {
		return errorResponse("Feature not enabled", http.StatusForbidden)
	}
	release, err := helmclient.DeleteRelease(params.ReleaseName)
	if err != nil {
		return errorResponse(fmt.Sprintf("Can't delete the release: %s", err), http.StatusBadRequest)
	}
	resource := makeReleaseResource(release.Release)
	payload := handlers.DataResourceBody(resource)
	return releasesapi.NewDeleteReleaseOK().WithPayload(payload)
}

func errorResponse(message string, errorCode int64) middleware.Responder {
	return releasesapi.NewGetAllReleasesDefault(int(errorCode)).WithPayload(
		&models.Error{Code: helpers.Int64ToPtr(errorCode), Message: &message},
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
		ChartIcon:    &release.Chart.Metadata.Icon,
		Updated:      helpers.StrToPtr(timeconv.String(release.Info.LastDeployed)),
		Name:         &release.Name,
		Namespace:    &release.Namespace,
		Status:       helpers.StrToPtr(release.Info.Status.Code.String()),
	}
	return &ret
}

func makeReleaseExtendedResource(release *hapi_release5.Release) *models.Resource {
	var ret models.Resource
	if release == nil {
		return &ret
	}
	ret.Type = helpers.StrToPtr("release")
	ret.ID = helpers.StrToPtr(release.Name)
	ret.Attributes = &models.ReleaseExtended{
		ChartName:    &release.Chart.Metadata.Name,
		ChartVersion: &release.Chart.Metadata.Version,
		ChartIcon:    &release.Chart.Metadata.Icon,
		Updated:      helpers.StrToPtr(timeconv.String(release.Info.LastDeployed)),
		Name:         &release.Name,
		Namespace:    &release.Namespace,
		Status:       helpers.StrToPtr(release.Info.Status.Code.String()),
		Resources:    helpers.StrToPtr(release.Info.Status.Resources),
		Notes:        helpers.StrToPtr(release.Info.Status.Notes),
	}
	return &ret
}
