package mocks

import (
	"errors"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/kubernetes-helm/monocular/src/api/data"
	releasesapi "github.com/kubernetes-helm/monocular/src/api/swagger/restapi/operations/releases"
	"k8s.io/helm/pkg/proto/hapi/chart"
	helm_releases "k8s.io/helm/pkg/proto/hapi/release"
	rls "k8s.io/helm/pkg/proto/hapi/services"
)

type mockedClient struct{}
type mockedBrokenClient mockedClient

// Resource represents a Helm release resource
var Resource = helm_releases.Release{
	Name:      "my-release-name",
	Namespace: "my-namespace",
	Chart: &chart.Chart{
		Metadata: &chart.Metadata{
			Name:    "my-chart",
			Version: "1.2.3",
			Icon: "chart-icon",
		},
	},
	Info: &helm_releases.Info{
		LastDeployed: &timestamp.Timestamp{},
		Status: &helm_releases.Status{
			Code:      200,
			Resources: "my-resources",
			Notes:     "my-notes",
		},
	},
}

// NewMockedClient returns a mocked version of the Helm Client
func NewMockedClient() data.Client {
	return &mockedClient{}
}

// NewMockedBrokenClient Fails to initialize
func NewMockedBrokenClient() data.Client {
	return &mockedBrokenClient{}
}

func (c *mockedClient) ListReleases(params releasesapi.GetAllReleasesParams) (*rls.ListReleasesResponse, error) {
	return &rls.ListReleasesResponse{
		Releases: []*helm_releases.Release{&Resource},
	}, nil
}

func (c *mockedClient) InstallRelease(chartPath string, params releasesapi.CreateReleaseParams) (*rls.InstallReleaseResponse, error) {
	return &rls.InstallReleaseResponse{Release: &Resource}, nil
}

func (c *mockedBrokenClient) ListReleases(params releasesapi.GetAllReleasesParams) (*rls.ListReleasesResponse, error) {
	return nil, errors.New("Can't initialize")
}

func (c *mockedBrokenClient) InstallRelease(chartPath string, params releasesapi.CreateReleaseParams) (*rls.InstallReleaseResponse, error) {
	return nil, errors.New("Can't initialize")
}

func (c *mockedClient) DeleteRelease(releaseName string) (*rls.UninstallReleaseResponse, error) {
	return &rls.UninstallReleaseResponse{Release: &Resource}, nil
}

func (c *mockedClient) GetRelease(releaseName string) (*rls.GetReleaseContentResponse, error) {
	return &rls.GetReleaseContentResponse{}, nil
}

func (c *mockedBrokenClient) DeleteRelease(releaseName string) (*rls.UninstallReleaseResponse, error) {
	return nil, errors.New("Can't initialize")
}

func (c *mockedBrokenClient) GetRelease(releaseName string) (*rls.GetReleaseContentResponse, error) {
	return nil, errors.New("Can't initialize")
}
