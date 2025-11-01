FROM golang:1-alpine3.22 AS builder
WORKDIR /go/src/entrypoint
COPY ./go.mod /go/src/entrypoint
COPY ./entrypoint.go /go/src/entrypoint/
RUN go build  -o entrypoint

FROM traefik:v1.7-alpine
COPY --from=builder /go/src/entrypoint/entrypoint /
COPY ./traefik.toml /etc/traefik/traefik.toml
RUN mkdir /cert
ENTRYPOINT [ "/entrypoint" ]
CMD ["/usr/local/bin/traefik"]
