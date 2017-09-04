[![Build Status](https://travis-ci.org/kubernetes-helm/monocular.svg?branch=master)](https://travis-ci.org/kubernetes-helm/monocular)
[![codecov](https://codecov.io/gh/kubernetes-helm/monocular/branch/master/graph/badge.svg)](https://codecov.io/gh/kubernetes-helm/monocular)
[![Go Report Card](https://goreportcard.com/badge/github.com/kubernetes-helm/monocular)](https://goreportcard.com/report/github.com/kubernetes-helm/monocular)
[![codebeat badge](https://codebeat.co/badges/2443005b-8e19-428a-8c67-14a2af60e7fd)](https://codebeat.co/projects/github-com-kubernetes-helm-monocular-master)

# Monocular API

The API is a golang HTTP RESTFul server. It abstracts away Helm Chart Repository data and provides a simple, idiomatic HTTP interface for search and discovery functionality. E.g.:

- search for official community "stable" and "incubator" charts
- get detailed version information on particular repo/charts
- browse charts in a repo

All commands and relative directories below assume a current working directory at the API source code root, i.e.:

- `$GOPATH/src/github.com/kubernetes-helm/monocular/src/api/`

# Building Monocular

`Makefile` provides a convenience for building locally:

- `make build`

The resulting will be placed inside `rootfs/usr/bin`, which is not coincidentally where `Dockerfile` assumes a `monocular` executable will be when building images.

# Building Docker Images

To build a docker image locally:

- `make docker-build`

The image will be tagged as `bitnami/monocular-api:latest` by default. Set `IMAGE_REPO` and `IMAGE_TAG` to override this.

# Running Monocular

To launch without building:
```
$ PORT=8080 go run main.go
serving monocular at http://127.0.0.1:8080
```

# Updating the API specification using swagger

Monocular uses [go-swagger](https://github.com/go-swagger/go-swagger) to define and generate the RESTFul server code. `Makefile` provides a convenience for generating server stub code:

- `make swagger-serverstub`

# Testing the API

- `make test`
