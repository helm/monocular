# Monocular

- Under ACTIVE development

Monocular is a search and discovery interface for Helm Charts Repositories. Charts Repositories are collections of curated Kubernetes application definitions.

A monocular is a single-lensed telescope, and perhaps a twee synonym for the kind of kit that would be useful when inspecting a disordered stack of nautical charts, like a magnifying glass or microscope with a spyglass aesthetic. Its adjectival form suggests a Cyclops, with whom Oddysseus was definitely on familiar terms; and its closest synonym is monocle, a well-worn accoutrement of Victorian Great Britain, surely the greatest of all naval empires. `kubernetes` indeed.

As with The Beatles, some things are better in mono.

Visit [the Helm repository](https://github.com/kubernetes/helm) to learn more about Helm, the package manager for Kubernetes.

Visit [the Charts repository](https://github.com/kubernetes/charts) to learn more about Charts, the Helm unit of configuration for a Kubernetes application definition.

# Prerequisites

Monocular is a golang RESTFul HTTP server.

`Makefile` assumes [docker](https://www.docker.com) for containerized development; and [glide](http://glide.sh) for dependency enforcement.

`make bootstrap` will launch a docker container, and `run` a `glide install` command to install all dependencies to the `vendor/` directory.

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

## Status of the Project

This project is still under active development, so you'll likely encounter [issues](https://github.com/helm/monocular/issues). Please participate by filing issues or contributing a pull request!
