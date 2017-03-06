package releases

import (
	log "github.com/Sirupsen/logrus"
	releasesapi "github.com/helm/monocular/src/api/swagger/restapi/operations/releases"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/proto/hapi/services"
	rls "k8s.io/helm/pkg/proto/hapi/services"
)

// ListReleases returns the list of helm releases
func ListReleases(client *helm.Client, params releasesapi.GetAllReleasesParams) (*rls.ListReleasesResponse, error) {
	stats := []release.Status_Code{
		release.Status_DEPLOYED,
	}
	resp, err := client.ListReleases(
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
func InstallRelease(client *helm.Client, chartPath string, params releasesapi.CreateReleaseParams) (*rls.InstallReleaseResponse, error) {
	ns := params.Data.Namespace
	if ns == "" {
		ns = "default"
	}

	return client.InstallRelease(
		chartPath,
		ns,
		helm.ValueOverrides([]byte{}),
		helm.ReleaseName(params.Data.ReleaseName),
		helm.InstallDryRun(params.Data.DryRun))
}
