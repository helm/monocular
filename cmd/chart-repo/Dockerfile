FROM golang:1.11 as builder
COPY . /go/src/github.com/helm/monocular
WORKDIR /go/src/github.com/helm/monocular

ARG CODE_VERSION
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags "-X main.version=$CODE_VERSION" ./cmd/chart-repo

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/helm/monocular/chart-repo /chart-repo
CMD ["/chart-repo"]
