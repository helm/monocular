# Monocular
[![CircleCI](https://circleci.com/gh/helm/monocular.svg?style=svg)](https://circleci.com/gh/helm/monocular)

Monocular is a web-based application that enables the search and discovery of
charts from multiple Helm Chart repositories. It is the codebase that powers the
[Helm Hub](https://github.com/helm/hub) project.

![Monocular Screenshot](docs/MonocularScreenshot.gif)

Click [here](docs/about.md) to learn more about Helm, Charts and Kubernetes.

## Install

You can use the chart in this repository to install Monocular in your cluster.

### Prerequisites
- [Helm and Tiller installed](https://github.com/helm/helm/blob/master/docs/quickstart.md)
- [Nginx Ingress controller](https://kubeapps.com/charts/stable/nginx-ingress)
  - Install with Helm: `helm install stable/nginx-ingress`
  - **Minikube/Kubeadm**: `helm install stable/nginx-ingress --set controller.hostNetwork=true`


```console
$ helm repo add monocular https://helm.github.io/monocular
$ helm install monocular/monocular
```

### Access Monocular

Use the Ingress endpoint to access your Monocular instance:

```console
# Wait for all pods to be running (this can take a few minutes)
$ kubectl get pods --watch

$ kubectl get ingress
NAME                        HOSTS     ADDRESS         PORTS     AGE
tailored-alpaca-monocular   *         192.168.64.30   80        11h
```

Visit the address specified in the Ingress object in your browser, e.g. http://192.168.64.30.

Read more on how to deploy Monocular [here](chart/monocular/README.md).

## Documentation

- [Configuration](chart/monocular/README.md#configuration)
- [Deployment](chart/monocular/README.md)
- [Development](docs/development.md)

## Looking for an in-cluster Application management UI?

To focus on the CNCF Helm Hub requirements, in-cluster features have been
removed from Monocular 1.0 and above. We believe that providing a good solution
for deploying and managing apps in-cluster is an orthogonal user experience to a
public search and discovery site. There is other tooling that can support this
usecase better (e.g. [Kubeapps](https://github.com/kubeapps/kubeapps) or [RedHat
Automation
Broker](https://blog.openshift.com/automation-broker-discovering-helm-charts/)).

[Monocular v0.7.3](https://github.com/helm/monocular/releases/tag/v0.7.3)
includes in-cluster features and can still be installed and used until your team
has migrated to another tool.

## Roadmap

The [Monocular roadmap is currently located in the wiki](https://github.com/helm/monocular/wiki/Roadmap).

## Contribute

This project is still under active development, so you'll likely encounter
[issues](https://github.com/helm/monocular/issues).

Interested in contributing? Check out the [documentation](CONTRIBUTING.md).

Also see [developer's guide](docs/development.md) for information on how to
build and test the code.
