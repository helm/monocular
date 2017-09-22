package client

import (
	"errors"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/kubernetes-helm/monocular/src/api/config"
	"github.com/kubernetes-helm/monocular/src/api/data"
	releasesapi "github.com/kubernetes-helm/monocular/src/api/swagger/restapi/operations/releases"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/helm/portforwarder"
	"k8s.io/helm/pkg/kube"
	rls "k8s.io/helm/pkg/proto/hapi/services"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"

	helmreleases "github.com/kubernetes-helm/monocular/src/api/data/helm/releases"
	"k8s.io/kubernetes/pkg/client/restclient"
)

const (
	tillerServiceName = "tiller-deploy"
	tillerPort        = 44134
)

type helmClient struct {
	portForward bool
}

// NewHelmClient returns the Helm implementation of data.Client
func NewHelmClient() data.Client {
	conf, _ := config.GetConfig()
	return &helmClient{
		portForward: conf.TillerPortForward,
	}
}

// InitializeClient returns a helm.client
func (c *helmClient) initialize() (*helm.Client, error) {
	// From helm.setupConnection
	conf, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	tillerNamespace := conf.TillerNamespace
	var tillerHost = fmt.Sprintf("%s.%s:%d", tillerServiceName, tillerNamespace, tillerPort)

	if c.portForward {
		config, kubeClient, err := getKubeClient("")
		if err != nil {
			return nil, err
		}

		tunnel, err := portforwarder.New(tillerNamespace, kubeClient, config)
		if err != nil {
			return nil, err
		}

		log.WithFields(log.Fields{"host": "localhost", "port": tunnel.Local}).Info("Helm Client created")
		tillerHost = fmt.Sprintf("localhost:%d", tunnel.Local)
	}

	client := helm.NewClient(helm.Host(tillerHost))
	// test connection
	if _, err := client.GetVersion(); err != nil {
		return nil, errors.New("failed to connect to Tiller, are you sure it is installed?")
	}

	return client, nil
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

func (c *helmClient) GetRelease(releaseName string) (*rls.GetReleaseContentResponse, error) {
	client, err := c.initialize()
	if err != nil {
		return nil, err
	}
	return helmreleases.GetRelease(client, releaseName)
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
