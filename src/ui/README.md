# Monocular UI

The UI is a web client for the [Monocular
API](https://github.com/helm/monocular/tree/master/src/api), which exposes an easy way to
navigate and search [Helm Charts](https://github.com/kubernetes/charts).

Regarding its functionality we can highlight:

* Listing of available and curated charts.
* Search mechanism.
* Chart information page, which includes instructions on how to use the
  chart, how to install it, etc.

## Developers

### Running Monocular UI

Monocular UI requires a running instance of the Monocular backend.

The easiest way to have a running multi-tier development environment is to use the the `docker-compose.yml` file placed at the project root directory.

Refer to [Running a development environment](src/README.md) for more details.

### Stack

The web application is based on the components listed below.

* [Angular 2](https://angular.io/)
* [@angular/cli](https://github.com/angular/angular-cli)
* Typescript
* Sass
* [Webpack](https://webpack.github.io/)
* Bootstrap

### Components tree

See below a representation of the implemented Angular components tree.

![components tree](https://cloud.githubusercontent.com/assets/24523/23182395/3ff0382a-f82d-11e6-9b64-2b8b0a9e45e9.png)
