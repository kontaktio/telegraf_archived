FROM golang:1.13-alpine as builder

ARG BUILD_BRANCH=develop

RUN apk --update --no-cache add dep git make

#ENV DEP_VERSION 0.5.0

RUN go get 'github.com/kontaktio/telegraf'
#RUN go get -u github.com/golang/dep/...

RUN mv $GOPATH/src/github.com/kontaktio $GOPATH/src/github.com/influxdata
WORKDIR $GOPATH/src/github.com/influxdata/telegraf
RUN git checkout "${BUILD_BRANCH}" && \
    dep ensure -vendor-only && \
    make go-install

FROM alpine:latest
COPY --from=builder /go/bin/* /usr/bin/
COPY --from=hairyhenderson/gomplate:alpine /bin/gomplate /bin/gomplate

RUN apk add --update --no-cache ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

RUN wget -O /etc/ssl/ca-bundle.pem https://curl.haxx.se/ca/cacert.pem

COPY tools/build/telegraf.conf.tpl /telegraf.conf.tpl
COPY tools/build/entrypoint.sh /entrypoint.sh

EXPOSE 8080

ENTRYPOINT "/entrypoint.sh"
