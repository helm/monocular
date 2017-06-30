## What is Kubernetes?

Kubernetes is an open source container cluster manager used to deploy, scale and operate applications across a number of host
computers. Kubernetes provides an API and primitives for managing those applications and their associated resources. Kubernetes
is sometimes abbreviated to “k8s” in documentation and tutorials.

Learn more about Kubernetes at https://kubernetes.io

## What is Helm?

Helm is a tool for managing applications that run in the Kubernetes cluster manager. Helm provides a set of operations that are
useful for managing applications, for example, inspect, install, upgrade and delete. Helm aims to provide a similar experience to
package managers such as [apt](https://wiki.debian.org/Apt) or [homebrew](https://brew.sh/), but for Kubernetes apps.

Learn more about Helm at https://helm.sh

## What is a chart?

A helm chart describes how to manage a specific application on Kubernetes. It consists of metadata that describes the application
plus the infrastructure needed to operate it in terms of the standard Kubernetes primitives. Each chart references one or more
(typically docker-compatible) container images that contain the application code to be run.

Learn more about charts at https://github.com/kubernetes/helm/blob/master/docs/charts.md

## What is Monocular?

Monocular is a part of the Helm project and aims to provide a way to search for and discover apps that have been packaged in Helm
Charts. Monocular includes a scanning back-end for indexing charts and their metadata and a simple user interface.

Other resources:

- [Project README](https://github.com/kubernetes-helm/monocular/blob/master/README.md)
- [Project Background](https://deis.com/blog/2017/building-a-helm-ui/)
- [Technical Overview](https://engineering.bitnami.com/2017/02/22/what-the-helm-is-monocular.html)

The Charts indexed by Monocular are from the official Kubernetes Helm Chart repository.

The process for contributing new Charts can be found at: https://github.com/kubernetes/charts#contributing-a-chart

*Coming soon: Point Monocular to a Helm Chart registry of your choice.*


## Who are the maintainers of Monocular and KubeApps.com?

Bitnami and Deis are the main committers to the Monocular project. Bitnami also sponsors the [KubeApps](https://kubeapps.com) website by providing hosting on GKE and contributing to website design.
