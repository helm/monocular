package releases

import (
	log "github.com/Sirupsen/logrus"
	"github.com/helm/monocular/src/api/data"
	releasesapi "github.com/helm/monocular/src/api/swagger/restapi/operations/releases"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/proto/hapi/services"
	rls "k8s.io/helm/pkg/proto/hapi/services"
)

type helmReleases struct {
	client *helm.Client
}

// NewHelmReleases returns the Helm implementation for the interface data.Releases
func NewHelmReleases(cl *helm.Client) data.Releases {
	return &helmReleases{
		client: cl,
	}
}

// ListReleases returns the list of helm releases
func (r *helmReleases) ListReleases() (*rls.ListReleasesResponse, error) {
	stats := []release.Status_Code{
		release.Status_DEPLOYED,
	}
	resp, err := r.client.ListReleases(
		helm.ReleaseListFilter(""),
		helm.ReleaseListSort(int32(services.ListSort_LAST_RELEASED)),
		helm.ReleaseListOrder(int32(services.ListSort_DESC)),
		helm.ReleaseListStatuses(stats),
	)

	if err != nil {
		log.WithError(err).Error("Can't retrieve the list of releases")
		return nil, err
	}

	return resp, err
}

// InstallRelease wraps helms client installReleae method
func (r *helmReleases) InstallRelease(chartPath string, params releasesapi.CreateReleaseParams) (*rls.InstallReleaseResponse, error) {
	ns := params.Data.Namespace
	if ns == "" {
		ns = "default"
	}

	return r.client.InstallRelease(
		chartPath,
		ns,
		helm.ValueOverrides([]byte{}),
		helm.ReleaseName(params.Data.ReleaseName),
		helm.InstallDryRun(params.Data.DryRun))
}
