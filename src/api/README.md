[![Build Status](https://travis-ci.org/helm/monocular.svg?branch=master)](https://travis-ci.org/helm/monocular)
[![codecov](https://codecov.io/gh/helm/monocular/branch/master/graph/badge.svg)](https://codecov.io/gh/helm/monocular)
[![Go Report Card](https://goreportcard.com/badge/github.com/helm/monocular)](https://goreportcard.com/report/github.com/helm/monocular)
[![codebeat badge](https://codebeat.co/badges/820a0d9f-5282-4a8e-b5f7-a27c217d9f0e)](https://codebeat.co/projects/github-com-helm-monocular)

# Monocular API

The API is a golang HTTP RESTFul server. It abstracts away Helm Chart Repository data and provides a simple, idiomatic HTTP interface for search and discovery functionality. E.g.:

- search for official community "stable" and "incubator" charts
- get detailed version information on particular repo/charts
- browse charts in a repo

All commands and relative directories below assume a current working directory at the API source code root, i.e.:

- `$GOPATH/src/github.com/helm/monocular/src/api/`

# Building Monocular

`Makefile` provides a convenience for building locally:

- `make build`

The resulting will be placed inside `rootfs/usr/bin`, which is not coincidentally where `Dockerfile` assumes a `monocular` executable will be when building images.

# Building Docker Images

To build a docker image locally:

- `IMAGE_PREFIX=superdev make docker-build`

Currently, you must provide an `IMAGE_PREFIX` to properly associate the resultant image with a registry (e.g., dockerhub) account. The image will be tagged with the current short git SHA (e.g., `c1c0e7f`) for an "immutable" reference, and a "mutable" tag of `canary` to reflect "latest".

And to push to a public registry, assuming the image has been built on your system previously following the example above:

- `IMAGE_PREFIX=superdev make docker-push`

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
