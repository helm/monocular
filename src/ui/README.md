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

Currently, the UI is relying on mocked data so an instance of Monocular
API is not necessary for the time being. Expect this to change in the
near future.

In order to run the Angular application we provide a `docker-compose.yml`. Just execute:

```
docker-compose up
```

Once the initialization is completed you should see a message like:

```
** NG Live Development Server is running on http://0.0.0.0:4200. **
```

Now just visit `http://{your-docker-machine-ip-address}:4200` and enjoy!

### Stack

The web application is based on the components listed below.

* [Angular 2](https://angular.io/)
* [angular-cli](https://github.com/angular/angular-cli)
* Typescript
* Sass
* [Webpack](https://webpack.github.io/)
* Bootstrap

### Components tree

See below a representation of the implemented Angular components tree.

![components
tree](https://cloud.githubusercontent.com/assets/24523/18405323/15b30e42-76a6-11e6-8d3b-c29794d2e06c.png)
