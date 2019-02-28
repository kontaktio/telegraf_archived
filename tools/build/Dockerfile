FROM golang:1.11-alpine

ARG SCRIPTS_SOURCE=test
ARG BUILD_BRANCH=develop

RUN apk update
RUN apk add git python py-pip make bash nodejs nodejs-npm

RUN pip install requests influxdb awscli

RUN npm config set unsafe-perm true
RUN npm i -g pm2@latest

COPY tools/build/generate_and_start_telegraf.sh /generate_and_start_telegraf.sh
RUN chmod +x /generate_and_start_telegraf.sh

RUN go get 'github.com/kontaktio/telegraf'
RUN go get -u github.com/golang/dep/...

RUN mv $GOPATH/src/github.com/kontaktio $GOPATH/src/github.com/influxdata
WORKDIR $GOPATH/src/github.com/influxdata/telegraf
RUN git checkout "${BUILD_BRANCH}"
RUN make

RUN touch /var/log/telegraf-config-gen.log
RUN mkdir /etc/telegraf

CMD["/bin/bash", "/generate_and_start_telegraf.sh", "${SCRIPTS_SOURCE}"]
