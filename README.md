# Monocular
[![Build
Status](https://travis-ci.org/helm/monocular.svg?branch=master)](https://travis-ci.org/helm/monocular)

Monocular is web-based UI for managing Kubernetes applications packaged as Helm
Charts. It allows you to search and discover available charts from multiple
repositories, and install them in your cluster with one click.

![Monocular Screenshot](docs/MonocularScreenshot.gif)

See Monocular in action at [KubeApps.com](https://kubeapps.com) or click [here](docs/about.md) to learn more about Helm, Charts and Kubernetes.

##### Video links
- [Screencast](https://www.youtube.com/watch?v=YoEbvDrI5ng)
- [Helm and Monocular Webinar](https://www.youtube.com/watch?v=u8kDkHgRbWQ)

## Install

You can use the chart in this repository to install Monocular in your cluster.

##### Prerequisites
- [Helm and Tiller installed](https://github.com/kubernetes/helm/blob/master/docs/quickstart.md)
- [Nginx Ingress controller](https://github.com/kubernetes/ingress)

```console
$ git clone https://github.com/helm/monocular.git
$ cd ./monocular
$ helm install ./deployment/monocular
```

Read more on how to deploy Monocular [here](docs/deployment.md).

## Documentation

- [Configuration](docs/configuration.md)
- [Deployment](docs/deployment.md)
- [Development](docs/development.md)

## Contribute

This project is still under active development, so you'll likely encounter
[issues](https://github.com/helm/monocular/issues).

Interested in contributing? Check out the [documentation](CONTRIBUTING.md).

Also see [developer's guide](docs/development.md) for information on how to
build and test the code.
