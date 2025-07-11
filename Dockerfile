FROM golang:1.23 AS builder

COPY ./ /build
RUN cd /build \
    && CGO_ENABLED=0 go build -o /mailcow-exporter /build/cmd/main.go \
    && rm -Rf /build

FROM alpine:3.21

RUN apk add --no-cache \
        openssl \
        ca-certificates \
        libc6-compat
COPY --from=builder /mailcow-exporter /usr/local/bin/mailcow-exporter

ENTRYPOINT [ "/usr/local/bin/mailcow-exporter" ]
