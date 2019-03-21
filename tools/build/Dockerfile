FROM golang:1.11-alpine as builder

ARG SCRIPTS_SOURCE=test
ARG BUILD_BRANCH=develop

RUN apk add git make curl

ENV DEP_VERSION 0.5.0
RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v${DEP_VERSION}/dep-linux-amd64 && chmod +x /usr/local/bin/dep

RUN go get 'github.com/kontaktio/telegraf'
RUN go get -u github.com/golang/dep/...

RUN mv $GOPATH/src/github.com/kontaktio $GOPATH/src/github.com/influxdata
WORKDIR $GOPATH/src/github.com/influxdata/telegraf
RUN git checkout "${BUILD_BRANCH}"
RUN dep ensure -vendor-only
RUN make go-install

FROM alpine:3.6
COPY --from=builder /go/bin/* /usr/bin/

RUN apk update
RUN apk add python py-pip bash nodejs nodejs-npm

RUN pip install requests influxdb awscli

RUN npm config set unsafe-perm true
RUN npm i -g pm2@latest

COPY tools/build/generate_and_start_telegraf.sh /generate_and_start_telegraf.sh
RUN chmod +x /generate_and_start_telegraf.sh

RUN touch /var/log/telegraf-config-gen.log
RUN mkdir /etc/telegraf

ENV SCRIPTS_SOURCE=$SCRIPTS_SOURCE;
ENTRYPOINT "/generate_and_start_telegraf.sh" "test"
