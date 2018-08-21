FROM golang:1.10.2-alpine

COPY generate_and_start_telegrafs.sh /generate_and_start_telegrafs.sh
RUN chmod +x /generate_and_start_telegrafs.sh
WORKDIR /
RUN apk update
RUN apk add git python py-pip make bash
RUN pip install requests influxdb
RUN go get 'github.com/kontaktio/telegraf'
RUN mv /go/src/github.com/kontaktio /go/src/github.com/influxdata
WORKDIR /go/src/github.com/influxdata/telegraf
RUN make
RUN touch /var/log/telegraf-config-gen.log
RUN mkdir /etc/telegraf

CMD /bin/bash /generate_and_start_telegrafs.sh
