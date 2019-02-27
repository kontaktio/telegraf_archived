FROM golang:1.11-alpine

ARG SCRIPTS_SOURCE=test

RUN apk update
RUN apk add git python py-pip make bash nodejs

RUN pip install requests influxdb awscli

RUN npm config set unsafe-perm true
RUN npm i -g pm2@latest

COPY tools/build/generate_and_start_telegraf.sh /generate_and_start_telegraf.sh
RUN chmod +x /generate_and_start_telegraf.sh

RUN mkdir -p /go/src/github.com/influxdata/telegraf
WORKDIR /go/src/github.com/influxdata/telegraf
COPY * .

RUN touch /var/log/telegraf-config-gen.log
RUN mkdir /etc/telegraf

CMD /bin/bash /generate_and_start_telegraf.sh $SCRIPTS_SOURCE
