# Monocular

[Monocular](https://github.com/helm/monocular) is a web-based application that
enables the search and discovery of charts from multiple Helm Chart
repositories. It is the codebase that powers the [Helm
Hub](https://github.com/helm/hub) project.

## TL;DR;

```console
$ helm repo add monocular https://helm.github.io/monocular
$ helm install monocular/monocular
```

## Introduction

This chart bootstraps a [Monocular](https://github.com/helm/monocular) deployment on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

### [Nginx Ingress controller](https://github.com/kubernetes/ingress)

To avoid issues with Cross-Origin Resource Sharing, the Monocular chart sets up an Ingress resource to serve the frontend and the API on the same domain. This is used to route requests made to `<host>:<port>/` to the frontend pods, and `<host>:<port>/api` to the backend pods.

## Installing the Chart

First, ensure you have added the Monocular chart repository:

```console
$ helm repo add monocular https://helm.github.io/monocular
```

To install the chart with the release name `my-release`:

```console
$ helm install --name my-release monocular/monocular
```

The command deploys Monocular on the Kubernetes cluster in the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```console
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

See the [values](values.yaml) for the full list of configurable values.

### Configuring chart repositories

You can configure the chart repositories you want to see in Monocular with the `sync.repos` value, for example:

```console
$ cat > custom-repos.yaml <<EOF
sync:
  repos:
    - name: stable
      url: https://kubernetes-charts.storage.googleapis.com
      schedule: "0 * * * *"
      successfulJobsHistoryLimit: 1
    - name: incubator
      url: https://kubernetes-charts-incubator.storage.googleapis.com
      schedule: "*/5 * * * *"
    - name: monocular
      url: https://helm.github.io/monocular
EOF

`schedule` and `successfulJobsHistoryLimit` are optional parameters. They default to `"0 * * * *"` and `3` respectively

$ helm install monocular/monocular -f custom-repos.yaml
```

### Serve Monocular on a single domain

You can configure the Ingress object with the hostnames you wish to serve Monocular on:

```console
$ cat > custom-domains.yaml <<EOF
ingress:
  hosts:
  - monocular.local
EOF

$ helm install monocular/monocular -f custom-domains.yaml
```

### Other configuration options

|          Value          |               Description                |                                     Default                                     |
| ----------------------- | ---------------------------------------- | ------------------------------------------------------------------------------- |
| `sync.nodeSelector`     | `{}`                                     | Node labels for pod assignment                                                  |
| `sync.tolerations`      | Tolerations for pod assignment           | `[]`                                                                            |
| `sync.affinity`         | Node/Pod affinities                      | `{}`                                                                            |
| `chartsvc.replicas`     | Number of replicas for API service       | `3`                                                                             |
| `chartsvc.nodeSelector` | `{}`                                     | Node labels for pod assignment                                                  |
| `chartsvc.tolerations`  | Tolerations for pod assignment           | `[]`                                                                            |
| `chartsvc.affinity`     | Node/Pod affinities                      | `{}`                                                                            |
| `ui.replicaCount`       | Number of replicas for UI service        | `2`                                                                             |
| `ui.googleAnalyticsId`  | Google Analytics ID                      | `UA-XXXXXX-X` (unset)                                                           |
| `ui.appName`            | Name to use in title bar and header      | `Monocular`                                                                     |
| `ui.nodeSelector`       | Node labels for pod assignment           | `{}`                                                                            |
| `ui.tolerations`        | Tolerations for pod assignment           | `[]`                                                                            |
| `ui.affinity`           | Node/Pod affinities                      | `{}`                                                                            |
| `ingress.enabled`       | If enabled, create an Ingress object     | `true`                                                                          |
| `ingress.annotations`   | Ingress annotations                      | `{ingress.kubernetes.io/rewrite-target: /, kubernetes.io/ingress.class: nginx}` |
| `ingress.tls`           | TLS configuration for the Ingress object | `nil`                                                                           |
| `global.mongoUrl`       | External MongoDB connection URL          | `nil`                                                                           |
