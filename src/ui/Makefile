IMAGE_REPO ?= bitnami/monocular-ui
IMAGE_TAG ?= latest

install:
	docker-compose run --rm ui npm install

test:
	docker-compose run --rm ui npm test

test-ci:
	docker-compose run --rm ui npm run test-ci

compile:
	docker-compose run --rm ui npm run compile

compile-aot:
	docker-compose run --rm ui npm run compile-aot

docker-build: compile-aot
	docker build --rm -t ${IMAGE_REPO}:${IMAGE_TAG} rootfs/
