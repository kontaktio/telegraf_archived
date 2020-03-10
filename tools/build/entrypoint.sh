#!/bin/sh
set -e

if [[ ! -z "${INTERNAL_CONFIG}" ]]; then
  gomplate -f /telegraf.conf.tpl -o /telegraf.conf
  /usr/bin/telegraf -config /telegraf.conf
else
  /usr/bin/telegraf -config /etc/telegraf/telegraf.conf
fi
