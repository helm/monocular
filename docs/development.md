# Developers Guide

This guide explains how to set up your environment for developing on Helm and Tiller.

## Prerequisites
* Docker 1.10 or later
* A Kubernetes cluster with Helm/Tiller installed
* Telepresence 0.75 or later
* kubectl 1.2 or later (optional)
* Go 1.11 or later
* Go dep https://golang.github.io/dep/
* Git

## Architecture

The UI is an Angular 2 application located in `frontend/`. This path is mounted
into the UI container. The server watches for file changes and automatically
rebuilds the application.

* [UI documentation](../frontend/README.md)

The backend is a small Go REST API service, `chartsvc`, and background CronJobs
to run the `chart-repo` sync command.

## Running Monocular

We develop Monocular in a Kubernetes environment, in order to make use of the
CronJobs for syncing chart repositories. Minikube can be used to run a local
single-node cluster for developing Monocular:

```
$ minikube start
$ minikube addons enable ingress
$ helm init --wait
$ helm dependency update
$ helm install --name dev --namespace monocular ./chart/monocular
```

After a few minutes, you will be able to visit the Monocular in your browser
using the Ingress address (typically IP of Minikube VM).

### Starting Monocular development server

Use Telepresence to replace the UI Pod in your cluster with the development
server:

```
$ docker build -t monocular_ui ./dev_env/ui
$ telepresence --swap-deployment dev-monocular-ui --namespace monocular --expose 4200:8080 --docker-run --rm -ti -v $(pwd)/frontend:/app monocular_ui bash
```

Inside the container's bash shell, run the following command to start the
development server:

```
$ ng serve --host 0.0.0.0 --public-host https://localhost
```

Once running, refresh the Ingress address in your browser and you will be
connected to the development server.

## Developing chartsvc

chartsvc is a Go REST API service. To build it, run the following commands:

```
$ dep ensure
$ make -C cmd/chartsvc docker-build
```

Use Telepresence to run the an instance of the chartsvc locally:

```
$ telepresence --swap-deployment dev-monocular-chartsvc --namespace monocular --expose 8080:8080 --docker-run --rm -ti quay.io/helmpack/chartsvc /chartsvc --mongo-user=root --mongo-url=dev-mongodb
```

Note that the chartsvc should be rebuilt for new changes to take effect.

## Developing chart-repo

chart-repo is a CLI tool that is used within CronJobs to sync against Helm chart
repositories. In development, it can be run standalone outside of the scheduled
CronJobs. To build it, run the following commands:

```
$ dep ensure
$ make -C cmd/chart-repo docker-build
```

Use Telepresence to run the image, passing in the repository name and URL as
arguments:

```
$ export MONGO_PASSWORD=$(kubectl get secret --namespace monocular dev-mongodb -o jsonpath="{.data.mongodb-root-password}" | base64 --decode)
$ telepresence --namespace monocular --docker-run -e MONGO_PASSWORD=$MONGO_PASSWORD --rm -ti quay.io/helmpack/chart-repo /chart-repo sync --mongo-user=root --mongo-url=dev-mongodb stable https://kubernetes-charts.storage.googleapis.com
```

Note that the chart-repo should be rebuilt for new changes to take effect.
