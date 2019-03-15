#!/bin/sh

mkdir /config_generator
aws s3 cp s3://kontakt-telegraf-config/build-$1 /config_generator --recursive

cd /config_generator/event_processor
npm install

source /config_generator/env

cd /etc/telegraf

echo "API_URL: $API_URL"

python /config_generator/telegraf_event_config_generate.py \
    --influxdb-url "$INFLUXDB_URL" \
    --influxdb-port $INFLUXDB_PORT \
    --influxdb-username "$INFLUXDB_USERNAME" \
    --influxdb-password "$INFLUXDB_PASSWORD" \

pm2 start -f /usr/bin/telegraf -- --config telegraf.events.conf

cd /config_generator/event_processor
pm2 start -f node -- index.js