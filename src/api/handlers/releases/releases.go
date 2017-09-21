package releases

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-openapi/strfmt"
	"github.com/kubernetes-helm/monocular/src/api/data"
	"github.com/kubernetes-helm/monocular/src/api/data/cache/charthelper"
	"github.com/kubernetes-helm/monocular/src/api/data/pointerto"
	"github.com/kubernetes-helm/monocular/src/api/handlers"
	"github.com/kubernetes-helm/monocular/src/api/handlers/renderer"
	"github.com/kubernetes-helm/monocular/src/api/swagger/models"
	releasesapi "github.com/kubernetes-helm/monocular/src/api/swagger/restapi/operations/releases"
	hapi_release5 "k8s.io/helm/pkg/proto/hapi/release"
	rls "k8s.io/helm/pkg/proto/hapi/services"
	"k8s.io/helm/pkg/timeconv"
)

// ReleaseHandlers defines handlers that serve Helm release data
type ReleaseHandlers struct {
	chartsImplementation data.Charts
	helmClient           data.Client
}

// NewReleaseHandlers takes a data.Client implementation and returns a ReleaseHandlers struct
func NewReleaseHandlers(ch data.Charts, hc data.Client) *ReleaseHandlers {
	return &ReleaseHandlers{helmClient: hc, chartsImplementation: ch}
}

// GetReleases returns all the existing releases in your cluster
func (r *ReleaseHandlers) GetReleases(w http.ResponseWriter, req *http.Request) {
	releases, err := r.helmClient.ListReleases(releasesapi.GetAllReleasesParams{})
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	resources := makeReleaseResources(releases)
	payload := handlers.DataResourcesBody(resources)
	renderer.Render.JSON(w, http.StatusOK, payload)
}

// GetRelease returns the extended version of a release
func (r *ReleaseHandlers) GetRelease(w http.ResponseWriter, req *http.Request, params handlers.Params) {
	release, err := r.helmClient.GetRelease(params["releaseName"])
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	resource := makeReleaseExtendedResource(release.Release)
	payload := handlers.DataResourceBody(resource)
	renderer.Render.JSON(w, http.StatusOK, payload)
}

// CreateRelease installs a chart version
func (r *ReleaseHandlers) CreateRelease(w http.ResponseWriter, req *http.Request) {
	// Params validation
	format := strfmt.NewFormats()
	var params releasesapi.CreateReleaseBody
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "unable to parse request body")
		return
	}
	err = params.Validate(format)
	if err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	idSplit := strings.Split(*params.ChartID, "/")
	if len(idSplit) != 2 || idSplit[0] == "" || idSplit[1] == "" {
		errorResponse(w, http.StatusBadRequest, "chartId must include the repository name. i.e: stable/wordpress")
		return
	}

	// Search chart package and get local path
	repo, chartName := idSplit[0], idSplit[1]
	chartPackage, err := r.chartsImplementation.ChartVersionFromRepo(repo, chartName, *params.ChartVersion)
	if err != nil {
		errorResponse(w, http.StatusNotFound, "404 chart not found")
		return
	}
	chartPath := charthelper.TarballPath(chartPackage)

	release, err := r.helmClient.InstallRelease(chartPath, releasesapi.CreateReleaseParams{Data: params})
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Can't create the release: "+err.Error())
		return
	}

	resource := makeReleaseResource(release.Release)
	payload := handlers.DataResourceBody(resource)
	renderer.Render.JSON(w, http.StatusCreated, payload)
}

// DeleteRelease deletes an existing release
func (r *ReleaseHandlers) DeleteRelease(w http.ResponseWriter, req *http.Request, params handlers.Params) {
	release, err := r.helmClient.DeleteRelease(params["releaseName"])
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Can't delete the release: "+err.Error())
		return
	}
	resource := makeReleaseResource(release.Release)
	payload := handlers.DataResourceBody(resource)
	renderer.Render.JSON(w, http.StatusOK, payload)
}

func errorResponse(w http.ResponseWriter, errorCode int64, message string) {
	renderer.Render.JSON(w, int(errorCode),
		models.Error{Code: pointerto.Int64(errorCode), Message: &message})
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
	ret.Type = pointerto.String("release")
	ret.ID = pointerto.String(release.Name)
	ret.Attributes = &models.Release{
		ChartName:    &release.Chart.Metadata.Name,
		ChartVersion: &release.Chart.Metadata.Version,
		ChartIcon:    &release.Chart.Metadata.Icon,
		Updated:      pointerto.String(timeconv.String(release.Info.LastDeployed)),
		Name:         &release.Name,
		Namespace:    &release.Namespace,
		Status:       pointerto.String(release.Info.Status.Code.String()),
	}
	return &ret
}

func makeReleaseExtendedResource(release *hapi_release5.Release) *models.Resource {
	var ret models.Resource
	if release == nil {
		return &ret
	}
	ret.Type = pointerto.String("release")
	ret.ID = pointerto.String(release.Name)
	ret.Attributes = &models.ReleaseExtended{
		ChartName:    &release.Chart.Metadata.Name,
		ChartVersion: &release.Chart.Metadata.Version,
		ChartIcon:    &release.Chart.Metadata.Icon,
		Updated:      pointerto.String(timeconv.String(release.Info.LastDeployed)),
		Name:         &release.Name,
		Namespace:    &release.Namespace,
		Status:       pointerto.String(release.Info.Status.Code.String()),
		Resources:    pointerto.String(release.Info.Status.Resources),
		Notes:        pointerto.String(release.Info.Status.Notes),
	}
	return &ret
}
