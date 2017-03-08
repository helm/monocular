package client

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/helm/monocular/src/api/data"
	releasesapi "github.com/helm/monocular/src/api/swagger/restapi/operations/releases"
	"k8s.io/helm/cmd/helm/installer"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/helm/portforwarder"
	"k8s.io/helm/pkg/kube"
	rls "k8s.io/helm/pkg/proto/hapi/services"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"

	helmreleases "github.com/helm/monocular/src/api/data/helm/releases"
	kerrors "k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/client/restclient"
)

const tillerNamespace = "kube-system"

type helmClient struct{}

// NewHelmClient returns the Helm implementation of data.Client
func NewHelmClient() data.Client {
	return &helmClient{}
}

// InitializeClient returns a helm.client
func (c *helmClient) initialize() (*helm.Client, error) {
	// From helm.setupConnection
	config, client, err := getKubeClient("")
	if err != nil {
		return nil, err
	}

	tunnel, err := portforwarder.New(tillerNamespace, client, config)
	if err != nil {
		// Initialize tiller if it does not exist
		if err = installer.Install(client, tillerNamespace, "", false, false); err != nil {
			if !kerrors.IsAlreadyExists(err) {
				return nil, fmt.Errorf("error installing Tiller: %s", err)
			}
		}
		log.WithFields(log.Fields{"error": err.Error()}).Info("Can't connect to Tiller. Installing")
		return nil, fmt.Errorf("Can't find a working Tiller pod. Installing, please try again later")
	}

	log.WithFields(log.Fields{"host": "localhost", "port": tunnel.Local}).Info("Helm Client created")

	tillerHost := fmt.Sprintf("localhost:%d", tunnel.Local)
	return helm.NewClient(helm.Host(tillerHost)), nil
}

func (c *helmClient) ListReleases(params releasesapi.GetAllReleasesParams) (*rls.ListReleasesResponse, error) {
	client, err := c.initialize()
	if err != nil {
		return nil, err
	}
	return helmreleases.ListReleases(client, params)
}

func (c *helmClient) InstallRelease(chartPath string, params releasesapi.CreateReleaseParams) (*rls.InstallReleaseResponse, error) {
	client, err := c.initialize()
	if err != nil {
		return nil, err
	}
	return helmreleases.InstallRelease(client, chartPath, params)
}

func (c *helmClient) DeleteRelease(releaseName string) (*rls.UninstallReleaseResponse, error) {
	client, err := c.initialize()
	if err != nil {
		return nil, err
	}
	return helmreleases.DeleteRelease(client, releaseName)
}

// getKubeClient is a convenience method for creating kubernetes config and client
// for a given kubeconfig context
// TODO, this is not needed once we are in the same cluster.
func getKubeClient(context string) (*restclient.Config, *internalclientset.Clientset, error) {
	config, err := kube.GetConfig(context).ClientConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("could not get kubernetes config for context '%s': %s", context, err)
	}
	client, err := internalclientset.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get kubernetes client: %s", err)
	}
	return config, client, nil
}
