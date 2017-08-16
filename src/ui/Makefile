IMAGE_REPO ?= bitnami/monocular-ui
IMAGE_TAG ?= latest
VERSION ?= dev

install:
	npm install

test:
	npm test

test-ci:
	npm run test-ci

compile:
	npm run compile

compile-aot:
	npm run compile-aot

docker-build: compile-aot
	docker build --pull --rm -t ${IMAGE_REPO}:${IMAGE_TAG} rootfs/

set-version:
	sed -i src/version.ts -e 's/dev/${VERSION}/'
