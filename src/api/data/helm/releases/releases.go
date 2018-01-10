package releases

import (
	log "github.com/Sirupsen/logrus"
	releasesapi "github.com/kubernetes-helm/monocular/src/api/swagger/restapi/operations/releases"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/helm/cmd/helm/strvals"
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

// GetRelease gets the information of an existing release
func GetRelease(client *helm.Client, releaseName string) (*rls.GetReleaseContentResponse, error) {
	// TODO, find a way to retrieve all the information in a single call
	// We get the information about the release
	release, err := client.ReleaseContent(releaseName)
	if err != nil {
		return nil, err
	}

	// Now we populate the resources string
	status, err := client.ReleaseStatus(releaseName)
	if err != nil {
		return nil, err
	}
	release.Release.Info = status.Info
	return release, err
}

// InstallRelease wraps helms client installReleae method
func InstallRelease(client *helm.Client, chartPath string, params releasesapi.CreateReleaseParams) (*rls.InstallReleaseResponse, error) {
	ns := params.Data.Namespace
	if ns == "" {
		ns = "default"
	}

	var overrides, err = MakeOverrides(params.Data.Values)
	if err != nil {
		return nil, err
	}

	var result, errRel = client.InstallRelease(
		chartPath,
		ns,
		helm.ValueOverrides(overrides),
		helm.ReleaseName(params.Data.ReleaseName),
		helm.InstallDryRun(params.Data.DryRun))

	return result, errRel
}

// make overrides
func MakeOverrides(values string) ([]byte, error) {
	if values == "" {
		return []byte{}, nil
	}
	var overridesMap, errMap = strvals.Parse(values)
	if errMap != nil {
		return nil, errMap
	}
	return yaml.Marshal(overridesMap)
}

// InstallRelease wraps helms client installReleae method
func UpdateRelease(client *helm.Client, rlsName string, chartPath string, params releasesapi.CreateReleaseParams) (*rls.UpdateReleaseResponse, error) {
	var overrides, err = MakeOverrides(params.Data.Values)
	if err != nil {
		return nil, err
	}
	return client.UpdateRelease(rlsName,
		chartPath,
		helm.UpdateValueOverrides(overrides),
		helm.UpgradeDryRun(params.Data.DryRun))
}

// DeleteRelease deletes an existing helm chart
func DeleteRelease(client *helm.Client, releaseName string) (*rls.UninstallReleaseResponse, error) {
	opts := []helm.DeleteOption{
		helm.DeleteDryRun(false),
		helm.DeletePurge(false),
		helm.DeleteTimeout(300),
	}
	return client.DeleteRelease(releaseName, opts...)
}
