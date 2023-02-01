FROM golang:1.19-alpine as builder

ARG SCRIPTS_SOURCE=test
ARG BUILD_BRANCH=develop-125

RUN apk --update upgrade && \
    apk add git make curl

COPY . $HOME/src/
WORKDIR $HOME/src
RUN make build_tools
RUN ./tools/custom_builder/custom_builder --config ./tools/build/example.conf


FROM alpine:3.9
COPY --from=builder $HOME/src/telegraf /usr/bin/

ARG SCRIPTS_SOURCE=test

RUN apk update
RUN apk add
RUN apk --update upgrade && \
    apk add python py-pip ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/* && \
    wget -O /etc/ssl/ca-bundle.pem https://curl.haxx.se/ca/cacert.pem


RUN pip install awscli

COPY tools/build/start_telegraf_and_acceptor.sh /start_telegraf_and_acceptor.sh
RUN chmod +x /start_telegraf_and_acceptor.sh
#
RUN mkdir /root/.aws
#
EXPOSE 8080
ENV SCRIPTS_SOURCE=$SCRIPTS_SOURCE
#
ENTRYPOINT "/start_telegraf_and_acceptor.sh" "${SCRIPTS_SOURCE}"
