#!/bin/sh

mkdir /config
aws s3 cp s3://kontakt-telegraf-config/build-$1/env /config
aws s3 cp s3://kontakt-telegraf-config/telegraf.eventprocessor.$1.conf /config/telegraf.eventprocesor.conf

source /config/env

echo "API_URL: $API_URL"

pm2 start -f /usr/bin/telegraf -- --config /config/telegraf.eventprocessor.conf

cd /config_generator/event_processor
pm2 start -f node -- index.js /tmp/telegraf.sock

tail -f /var/log/telegraf-config-gen.log