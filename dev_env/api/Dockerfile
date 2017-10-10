FROM quay.io/deis/go-dev:v1.5.0

RUN go get github.com/codegangsta/gin

COPY rootfs/ /

WORKDIR /go/src/github.com/kubernetes-helm/monocular/src/api

ENTRYPOINT ["/app-entrypoint.sh"]

CMD ["gin", "-i", "run", "main.go"]
