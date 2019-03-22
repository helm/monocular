FROM golang:1.12 as builder
COPY . /go/src/github.com/helm/monocular
WORKDIR /go/src/github.com/helm/monocular

ARG VERSION
RUN GO111MODULE=on GOPROXY=https://gocenter.io CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags "-X main.version=$VERSION" ./cmd/chart-repo

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/helm/monocular/chart-repo /chart-repo
USER 1001
CMD ["/chart-repo"]
