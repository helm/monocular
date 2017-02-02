MUTABLE_VERSION ?= canary
VERSION ?= git-$(shell git rev-parse --short HEAD)
IMAGE_REGISTRY ?= # default to dockerhub
IMAGE_PREFIX ?= # we rely upon the user providing an IMAGE_PREFIX, e.g., IMAGE_PREFIX=jackfrancis make docker-push
IMAGE_NAME ?= ${SHORT_NAME}-api

IMAGE := ${IMAGE_REGISTRY}${IMAGE_PREFIX}/${IMAGE_NAME}:${VERSION}
MUTABLE_IMAGE := ${IMAGE_REGISTRY}${IMAGE_PREFIX}/${IMAGE_NAME}:${MUTABLE_VERSION}

.PHONY: docker-push
docker-push: docker-immutable-push docker-mutable-push

.PHONY: docker-immutable-push
docker-immutable-push:
	docker push ${IMAGE}

.PHONY: docker-mutable-push
docker-mutable-push:
	docker push ${MUTABLE_IMAGE}
