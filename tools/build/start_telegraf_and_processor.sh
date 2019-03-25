#!/bin/sh

mkdir /config_generator
aws s3 cp s3://kontakt-telegraf-config/build-$1 /config_generator --recursive
aws s3 cp s3://kontakt-telegraf-config/telegraf.eventprocessor.$1.conf /telegraf.eventprocesor.conf

source /config_generator/env

echo "API_URL: $API_URL"

pm2 start -f /usr/bin/telegraf -- --config /telegraf.eventprocessor.conf

cd /config_generator/event_processor
pm2 start -f node -- index.js /tmp/telegraf.sock

tail -f /var/log/telegraf-config-gen.log