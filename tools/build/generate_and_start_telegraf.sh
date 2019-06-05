#!/bin/sh

mkdir /config_generator
aws s3 cp s3://kontakt-telegraf-config/build-$1 /config_generator --recursive
source /config_generator/env

cd /etc/telegraf

echo "API_URL: $API_URL"

python /config_generator/telegraf_stream_config_generate.py \
    --api-key "$API_KEY" --api-url "$API_URL" \
    --influxdb-url "$INFLUXDB_URL" \
    --influxdb-port $INFLUXDB_PORT \
    --influxdb-username "$INFLUXDB_USERNAME" \
    --influxdb-password "$INFLUXDB_PASSWORD"  \
    --mqtt-url "$MQTT_URL" \
    --telemetry-types "$TELEMETRY_TYPES" \
    --streams-per-telegraf $STREAMS_PER_TELEGRAF \
    --api-venue-id "$VENUE_ID"
for f in /etc/telegraf/telegraf.stream.conf.*;
do
	pm2 start -f /usr/bin/telegraf -- --config $f
done;

python /config_generator/telegraf_location_config_generate.py \
    --api-key "$API_KEY" --api-url "$API_URL" \
    --influxdb-url "$INFLUXDB_URL" \
    --influxdb-port $INFLUXDB_PORT \
    --influxdb-username "$INFLUXDB_USERNAME" \
    --influxdb-password "$INFLUXDB_PASSWORD"  \
    --api-venue-id "$VENUE_ID"
for f in /etc/telegraf/telegraf.location.conf.*;
do
	pm2 start -f /usr/bin/telegraf -- --config $f
done;


python /config_generator/kapacitor_reports_job.py \
    --api-key "$API_KEY" --api-url "$API_URL" \
    --kapacitor-url "$KAPACITOR_URL" \
    --kapacitor-user "$KAPACITOR_USER" \
    --kapacitor-pass "$KAPACITOR_PASS"


tail -f /var/log/telegraf-config-gen.log