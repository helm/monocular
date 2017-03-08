package client

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/helm/monocular/src/api/config"
	"github.com/helm/monocular/src/api/data"
	releasesapi "github.com/helm/monocular/src/api/swagger/restapi/operations/releases"
	"k8s.io/helm/cmd/helm/installer"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/helm/portforwarder"
	"k8s.io/helm/pkg/kube"
	rls "k8s.io/helm/pkg/proto/hapi/services"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/typed/core/internalversion"
	"k8s.io/kubernetes/pkg/labels"

	helmreleases "github.com/helm/monocular/src/api/data/helm/releases"
	kerrors "k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/client/restclient"
)

const (
	tillerNamespace   = "kube-system"
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

func (c *helmClient) ensureTillerInstalled(namespace string, client *internalclientset.Clientset) error {
	_, err := getFirstRunningPod(client, namespace, labels.Set{"app": "helm", "name": "tiller"}.AsSelector())
	if err == nil {
		return nil
	}

	log.WithFields(log.Fields{"error": err.Error()}).Info("Can't connect to Tiller. Installing")
	if err = installer.Install(client, namespace, "", false, false); err != nil {
		if !kerrors.IsAlreadyExists(err) {
			return fmt.Errorf("error installing Tiller: %s", err)
		}
	}
	return fmt.Errorf("Can't connect to Tiller. Installing, please wait")
}

// InitializeClient returns a helm.client
func (c *helmClient) initialize() (*helm.Client, error) {
	// From helm.setupConnection
	config, client, err := getKubeClient("")
	if err != nil {
		return nil, err
	}

	err = c.ensureTillerInstalled(tillerNamespace, client)
	if err != nil {
		return nil, err
	}

	var tillerHost = fmt.Sprintf("%s.%s:%d", tillerServiceName, tillerNamespace, tillerPort)

	if c.portForward {
		tunnel, err := portforwarder.New(tillerNamespace, client, config)
		if err != nil {
			return nil, err
		}

		log.WithFields(log.Fields{"host": "localhost", "port": tunnel.Local}).Info("Helm Client created")
		tillerHost = fmt.Sprintf("localhost:%d", tunnel.Local)
	}

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

func getFirstRunningPod(client internalversion.PodsGetter, namespace string, selector labels.Selector) (*api.Pod, error) {
	options := api.ListOptions{LabelSelector: selector}
	pods, err := client.Pods(namespace).List(options)
	if err != nil {
		return nil, err
	}
	if len(pods.Items) < 1 {
		return nil, fmt.Errorf("could not find tiller")
	}
	for _, p := range pods.Items {
		if api.IsPodReady(&p) {
			return &p, nil
		}
	}
	return nil, fmt.Errorf("could not find a ready tiller pod")
}
