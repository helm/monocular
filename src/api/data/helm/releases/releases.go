package releases

import (
	log "github.com/Sirupsen/logrus"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/proto/hapi/services"
	rls "k8s.io/helm/pkg/proto/hapi/services"
)

// ListReleases returns the list of helm releases
func ListReleases(client *helm.Client) (*rls.ListReleasesResponse, error) {
	stats := []release.Status_Code{
		release.Status_DEPLOYED,
	}
	resp, err := client.ListReleases(
		helm.ReleaseListLimit(5),
		helm.ReleaseListOffset(""),
		helm.ReleaseListFilter(""),
		helm.ReleaseListSort(int32(services.ListSort_NAME)),
		helm.ReleaseListOrder(int32(services.ListSort_ASC)),
		helm.ReleaseListStatuses(stats),
	)

	if err != nil {
		log.WithError(err).Error("Can't retrieve the list of releases")
		return nil, err
	}

	return resp, err
}
