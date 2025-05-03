FROM golang:1.22.4-alpine3.20 AS build

COPY . /go/src/github.com/andreimarcu/linx-server
WORKDIR /go/src/github.com/andreimarcu/linx-server

RUN set -ex \
        && apk add --no-cache --virtual .build-deps git \
        && go mod download \
        && go build -o /go/bin/linx-server . \
        && apk del .build-deps

FROM alpine:3.20

COPY --from=build /go/bin/linx-server /usr/local/bin/linx-server

ENV GOPATH /go
ENV SSL_CERT_FILE /etc/ssl/cert.pem

COPY static /go/src/github.com/andreimarcu/linx-server/static/
COPY templates /go/src/github.com/andreimarcu/linx-server/templates/

RUN mkdir -p /data/files && mkdir -p /data/meta && mkdir -p /data/locks && mkdir -p /data/custom_pages && chown -R 65534:65534 /data

VOLUME ["/data/files", "/data/meta", "/data/locks", "/data/custom_pages"]

EXPOSE 8080
USER nobody
ENTRYPOINT ["/usr/local/bin/linx-server", "-bind=0.0.0.0:8080", "-filespath=/data/files/", "-metapath=/data/meta/", "-lockspath=/data/locks/", "-custompagespath=/data/custom_pages"]
CMD ["-sitename=linx", "-allowhotlink"]
