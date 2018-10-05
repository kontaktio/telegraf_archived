#!/bin/sh

mkdir /config_generator
aws s3 cp s3://kontakt-bigdogs-cfg/telegraf-build-test /config_generator --recursive

cd /etc/telegraf
python /config_generator/telegraf_stream_config_generate.py --api-key $API_KEY --influxdb-url ${INFLUXDB_URL:=http://bdinflux.kontakt.io} --influxdb-port ${INFLUXDB_PORT:=8086} --influxdb-username $INFLUXDB_USERNAME --influxdb-password $INFLUXDB_PASSWORD --api-url ${API_URL:=http://bdapi.kontakt.io} --data-collection-interval ${STREAM_DATA_COLLECTION_INTERVAL:=30s} --flush-interval ${STREAM_FLUSH_INTERVAL:=10s} --flush-jitter ${STREAM_FLUSH_JITTER:=0s} --debug ${DEBUG_ENABLED:=true} --log-file ${LOG_FILE_NAME:=/var/log/telegraf-config-gen.log} --mqtt-url ${MQTT_URL:=ssl://bdrtls.kontakt.io:8083} --telemetry-types ${TELEMETRY_TYPES:=all} --streams-per-telegraf ${STREAMS_PER_TELEGRAF:=100} $(if [ ! -z $VENUE_ID ]; then echo "--api-venue-id $VENUE_ID"; else echo ""; fi;)
for f in /etc/telegraf/telegraf.stream.conf.*;
do
	pm2 start /go/src/github.com/influxdata/telegraf/telegraf -- --config $f
done;

python /config_generator/telegraf_location_config_generate.py --api-key $API_KEY --kapacitor-url ${KAPACITOR_URL:=http://bdinflux.kontakt.io:9090} --kapacitor-user ${KAPACITOR_USER:=kontaktio} --kapacitor-pass ${KAPACITOR_PASS:=notthepassword} --influxdb-url ${INFLUXDB_URL:=http://bdinflux.kontakt.io} --influxdb-port ${INFLUXDB_PORT:=8086} --influxdb-username $INFLUXDB_USERNAME --influxdb-password $INFLUXDB_PASSWORD --api-url ${API_URL:=http://bdapi.kontakt.io} --data-collection-interval ${INFSOFT_DATA_COLLECTION_INTERVAL:=5s} --flush-interval ${INFSOFT_FLUSH_INTERVAL:=10s} --flush-jitter ${INFSOFT_FLUSH_JITTER:=2s} --debug ${DEBUG_ENABLED:=true} --log-file ${LOG_FILE_NAME:=/var/log/telegraf-config-gen.log} --tx-power ${TX_POWER:=-77} $(if [ ! -z $VENUE_ID ]; then echo "--api-venue-id $VENUE_ID"; else echo ""; fi;)
for f in /etc/telegraf/telegraf.location.conf.*;
do
	pm2 start /go/src/github.com/influxdata/telegraf/telegraf -- --config $f
done;

tail -f /var/log/telegraf-config-gen.log