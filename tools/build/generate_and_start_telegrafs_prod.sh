#!/bin/sh

mkdir /config_generator
aws s3 cp s3://kontakt-telegraf-config/build-prod /config_generator --recursive

cd /etc/telegraf
python /config_generator/telegraf_stream_config_generate.py \
    --api-key $API_KEY --api-url ${API_URL:=http://api.kontakt.io} \
    --influxdb-url ${INFLUXDB_URL:=http://influx.kontakt.io} \
    --influxdb-port ${INFLUXDB_PORT:=8086} \
    --influxdb-username $INFLUXDB_USERNAME \
    --influxdb-password $INFLUXDB_PASSWORD  \
    --mqtt-url ${MQTT_URL:=ssl://ovs.kontakt.io:8083} \
    --telemetry-types ${TELEMETRY_TYPES:=all} \
    --streams-per-telegraf ${STREAMS_PER_TELEGRAF:=250} \
    $(if [ ! -z $VENUE_ID ]; then echo "--api-venue-id $VENUE_ID"; else echo ""; fi;)
for f in /etc/telegraf/telegraf.stream.conf.*;
do
	pm2 start -f /go/src/github.com/influxdata/telegraf/telegraf -- --config $f
done;

python /config_generator/telegraf_location_config_generate.py \
    --api-key $API_KEY --api-url ${API_URL:=http://api.kontakt.io}  \
    --kapacitor-url ${KAPACITOR_URL:=http://influx.kontakt.io:9090} \
    --kapacitor-user ${KAPACITOR_USER:=kontaktio} \
    --kapacitor-pass ${KAPACITOR_PASS:=notthepassword} \
    --influxdb-url ${INFLUXDB_URL:=http://influx.kontakt.io} \
    --influxdb-port ${INFLUXDB_PORT:=8086} \
    --influxdb-username $INFLUXDB_USERNAME \
    --influxdb-password $INFLUXDB_PASSWORD \
    --tx-power ${TX_POWER:=-77} $(if [ ! -z $VENUE_ID ]; then echo "--api-venue-id $VENUE_ID"; else echo ""; fi;)
for f in /etc/telegraf/telegraf.location.conf.*;
do
	pm2 start -f /go/src/github.com/influxdata/telegraf/telegraf -- --config $f
done;


python /config_generator/kapacitor_reports_job.py \
    --api-key $API_KEY --api-url ${API_URL:=http://api.kontakt.io}  \
    --kapacitor-url ${KAPACITOR_URL:=http://influx.kontakt.io:9090} \
    --kapacitor-user ${KAPACITOR_USER:=kontaktio} \
    --kapacitor-pass ${KAPACITOR_PASS:=notthepassword}


tail -f /var/log/telegraf-config-gen.log