FROM quay.io/deis/go-dev:v0.22.0

MAINTAINER Miguel Martinez <miguel@bitnami.com>

RUN curl https://glide.sh/get | sh

COPY rootfs/ /

WORKDIR /go/src/github.com/helm/monocular/src/api

ENV HOST=0.0.0.0
ENV PORT=8080

ENTRYPOINT ["/app-entrypoint.sh"]

CMD ["go", "run", "main.go"]
