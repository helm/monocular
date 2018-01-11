package data

import releasesapi "github.com/kubernetes-helm/monocular/src/api/swagger/restapi/operations/releases"
import rls "k8s.io/helm/pkg/proto/hapi/services"

// Client is an interface for managing Helm Chart releases
type Client interface {
	ListReleases(params releasesapi.GetAllReleasesParams) (*rls.ListReleasesResponse, error)
	InstallRelease(chartPath string, params releasesapi.CreateReleaseParams) (*rls.InstallReleaseResponse, error)
	UpdateRelease(rlsName string, chartPath string, params releasesapi.CreateReleaseParams) (*rls.UpdateReleaseResponse, error)
	DeleteRelease(releaseName string, purge bool) (*rls.UninstallReleaseResponse, error)
	GetRelease(releaseName string) (*rls.GetReleaseContentResponse, error)
}
