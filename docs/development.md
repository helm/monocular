# Developers Guide

This guide explains how to set up your environment for developing on Helm and Tiller.

## Prerequisites
* Docker 1.10 or later
* A Kubernetes cluster with Helm/Tiller installed
* Telepresence 0.75 or later
* kubectl 1.2 or later (optional)
* Git

## Running Monocular

We develop Monocular in a Kubernetes environment. Minikube can be used to run a
local single-node cluster for developing Monocular:

```
$ minikube start
$ minikube addons enable ingress
$ helm init --wait
$ helm install --name monocular ./chart/monocular
```

After a few minutes, you will be able to visit the Monocular in your browser
using the Ingress address (typically IP of Minikube VM).

### Starting Monocular development server

Use Telepresence to replace the UI Pod in your cluster with the development
server:

```
$ docker build -t monocular_ui ./dev_env/ui
$ telepresence --swap-deployment m-monocular-ui --namespace monocular --expose 4200:8080 --docker-run --rm -ti -v $(pwd)/frontend:/app monocular_ui bash
```

Inside the container's bash shell, run the following command to start the
development server:

```
$ ng serve --host 0.0.0.0 --public-host https://localhost
```

Once running, refresh the Ingress address in your browser and you will be
connected to the development server.

**TODO:** document development instructions for chartsvc and chart-repo.

## Architecture

The UI is an Angular 2 application located in `frontend/`. This path is mounted
into the UI container. The server watches for file changes and automatically
rebuilds the application.

* [UI documentation](../frontend/README.md)

The backend is a small Go REST API service, `chartsvc`, and background CronJobs
to run the `chart-repo` sync command.
