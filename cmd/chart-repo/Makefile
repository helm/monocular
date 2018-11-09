IMAGE_REPO ?= quay.io/helmpack/chart-repo
IMAGE_TAG ?= latest
# Version of the binary to be produced
VERSION ?= $$(git rev-parse HEAD)

docker-build:
	# We use the context of the root dir
	docker build --pull --rm -t ${IMAGE_REPO}:${IMAGE_TAG} --build-arg "VERSION=${VERSION}" -f Dockerfile ../../
