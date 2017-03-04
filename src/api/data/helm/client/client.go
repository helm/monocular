package client

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/helm/portforwarder"
	"k8s.io/helm/pkg/kube"
	"k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/client/restclient"
)

const tillerNamespace = "kube-system"

// CreateTillerClient returns a helm.client
func CreateTillerClient() (*helm.Client, error) {
	// From helm.setupConnection
	config, client, err := getKubeClient("")
	if err != nil {
		return nil, err
	}

	tunnel, err := portforwarder.New(tillerNamespace, client, config)
	if err != nil {
		return nil, err
	}

	tillerHost := fmt.Sprintf("localhost:%d", tunnel.Local)

	log.WithFields(log.Fields{
		"host": "localhost",
		"port": tunnel.Local,
	}).Info("Helm Client created")

	return helm.NewClient(helm.Host(tillerHost)), nil
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
