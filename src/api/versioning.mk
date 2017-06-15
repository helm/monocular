IMAGE_REPO ?= bitnami/${SHORT_NAME}-api
IMAGE_TAG ?= latest

IMAGE := ${IMAGE_REPO}:${IMAGE_TAG}

.PHONY: docker-push
docker-push:
	docker push ${IMAGE}
