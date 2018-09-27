# Monocular UI

The UI is a web client for the
[chartsvc](https://github.com/helm/monocular/tree/master/cmd/chartsvc), which
exposes an easy way to navigate and search [Helm
Charts](https://github.com/kubernetes/charts).

Features of the UI includes:

* Listing of available charts from multiple repositories.
* Search charts by name, keywords, maintainer, etc.
* View chart information, e.g. installation notes, usage, versions.

## Developers

### Running Monocular UI

Monocular UI requires a running instance of the `chartsvc` backend.

The easiest way to have a running multi-tier development environment is to
install the Helm chart in your development Kubernetes cluster.

Refer to [the Developer Guide](../docs/development.md) for more details.

### Stack

The web application is based on the components listed below.

* [Angular 2](https://angular.io/)
* [angular/cli](https://github.com/angular/angular-cli)
* Typescript
* Sass
* [Webpack](https://webpack.github.io/)
* Bootstrap

### Building

`Makefile` provides a convenience for building locally:

- `make compile-aot`

The resulting compiled static Angular application will be placed inside `rootfs/dist`, which is coincidentally where `rootfs/Dockerfile` expects to find it.

### Building Docker Images

To build a docker image locally:

- `make docker-build`

The image will be tagged as `bitnami/monocular-ui:latest` by default. Set `IMAGE_REPO` and `IMAGE_TAG` to override this.

### Components tree

See below a representation of the implemented Angular components tree.

![components tree](https://cloud.githubusercontent.com/assets/24523/23182395/3ff0382a-f82d-11e6-9b64-2b8b0a9e45e9.png)
