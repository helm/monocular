FROM quay.io/deis/go-dev:v0.22.0

MAINTAINER Miguel Martinez <miguel@bitnami.com>

RUN curl https://glide.sh/get | sh
RUN go get github.com/codegangsta/gin

COPY rootfs/ /

WORKDIR /go/src/github.com/kubernetes-helm/monocular/src/api

ENTRYPOINT ["/app-entrypoint.sh"]

CMD ["gin", "-i", "run", "main.go"]
