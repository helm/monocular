FROM golang:1.11 as builder
COPY . /go/src/github.com/helm/monocular
WORKDIR /go/src/github.com/helm/monocular
RUN CGO_ENABLED=0 go build -a -installsuffix cgo ./cmd/chartsvc

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/helm/monocular/chartsvc /chartsvc
EXPOSE 8080
CMD ["/chartsvc"]
