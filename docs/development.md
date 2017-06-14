# Development

## App Structure

Monocular comprises a UI front end, and a RESTFul HTTP back end API.

### UI Prerequisites

The UI is an angular 2 client application located in `src/ui/`.

More UI docs are [here](src/ui/README.md).

### API Prerequisites

The API is a golang HTTP server located in `src/api/`.

`Makefile` assumes [docker](https://www.docker.com) for containerized development; and [glide](http://glide.sh) for dependency enforcement.

`cd src/api/ && make bootstrap` will launch a docker container, and run a `glide install` command to install all API dependencies in the `src/api/vendor/` directory.

More API docs are [here](src/api/README.md).

## Running a development environment

We leverage [docker](https://www.docker.com) (via `docker-compose`) to provide a multi-tier setup for development.

Running `docker-compose up` from the root directory will expose:

* API backend endpoint via `http://{your-docker-machine-ip-address}:8080`
* UI frontend via `http://{your-docker-machine-ip-address}:4200`  

**IMPORTANT**:
* If your Docker Machine hostname is different than *localhost*, you need to change
the `backendHostname` value in the file `src/ui/src/app/shared/services/config.service.ts`.

You can restart individual services doing `docker-compose restart api|ui`
