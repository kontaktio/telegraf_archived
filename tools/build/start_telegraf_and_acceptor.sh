#!/bin/sh

mkdir /config_generator

if [[ "${CONFIG_FROM_S3}" == "true" ]]; then
    mkdir $(dirname "${TELEGRAF_CONFIG_PATH}")
    aws s3 cp "${configPath}" "${TELEGRAF_CONFIG_PATH}"
fi

/usr/bin/telegraf -config "${TELEGRAF_CONFIG_PATH}" --pprof-addr :8088