SHORT_NAME ?= monocular

include versioning.mk
include includes.mk

# dockerized development environment variables
REPO_PATH := github.com/kubernetes-helm/${SHORT_NAME}
DEV_ENV_IMAGE := quay.io/deis/go-dev:v1.4.0
SWAGGER_IMAGE := quay.io/goswagger/swagger:0.6.0
DEV_ENV_WORK_DIR := /go/src/${REPO_PATH}/src/api
DEV_ENV_PREFIX := docker run --rm -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR} --net=host
DEV_ENV_CMD := ${DEV_ENV_PREFIX} ${DEV_ENV_IMAGE}
SWAGGER_CMD := docker run --rm -e GOPATH=/go -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR} ${SWAGGER_IMAGE}
SWAGGER_GEN_DIR := swagger
SWAGGER_SPEC := swagger.yml
BINARY_DEST_DIR := rootfs/usr/bin
# Common flags passed into Go's linker.
LDFLAGS := "-s -X main.version=${VERSION}"
# Docker Root FS
BINDIR := ./rootfs

all:
	@echo "Use a Makefile to control top-level building of the project."

bootstrap:
	${DEV_ENV_CMD} dep ensure

update:
	${DEV_ENV_CMD} dep ensure -update

# This illustrates a two-stage Docker build. docker-compile runs inside of
# the Docker environment. Other alternatives are cross-compiling, doing
# the build as a `docker build`.
build:
	${DEV_ENV_PREFIX} -e CGO_ENABLED=0 ${DEV_ENV_IMAGE} go build -a -installsuffix cgo -ldflags ${LDFLAGS} -o ${BINARY_DEST_DIR}/${SHORT_NAME} *.go || exit 1
	@$(call check-static-binary,$(BINARY_DEST_DIR)/${SHORT_NAME})

swagger-serverstub:
	${SWAGGER_CMD} generate server -A ${SHORT_NAME} -t ${SWAGGER_GEN_DIR} -f ${SWAGGER_SPEC}
	mv ${SWAGGER_GEN_DIR}/cmd/${SHORT_NAME}-server/main.go .
	rm -Rf ${SWAGGER_GEN_DIR}/cmd

test:
	${DEV_ENV_CMD} sh -c 'go test -v ./...'

test-cover:
	${DEV_ENV_CMD} _scripts/test-cover.sh

test-native:
	go test -tags=testonly -cover ./...

docker-build: build
	docker build --rm -t ${IMAGE} rootfs
