#!/bin/sh
cd /etc/telegraf
python /go/src/github.com/influxdata/telegraf/tools/telegraf_config_generate.py --api-key $API_KEY --influxdb-url ${INFLUXDB_URL:=http://testinflux.kontakt.io} --influxdb-port ${INFLUXDB_PORT:=8086} --influxdb-username $INFLUXDB_USERNAME --influxdb-password $INFLUXDB_PASSWORD --api-url ${API_URL:=http://testapi.kontakt.io} --data-collection-interval ${DATA_COLLECTION_INTERVAL:=30} --debug ${DEBUG_ENABLED:=true} --log-file ${LOG_FILE_NAME:=/var/log/telegraf-config-gen.log} --mqtt-url ${MQTT_URL:=ssl://testrtls.kontakt.io:8083} --telemetry-types ${TELEMETRY_TYPES:=health sensor location accelerometer button} --streams-per-telegraf ${STREAMS_PER_TELEGRAF:=100} $(if [ ! -z $VENUE_ID ]; then echo "--api-venue-id $VENUE_ID"; else echo ""; fi;)
for f in /etc/telegraf/telegraf.conf.*;
do
	/go/src/github.com/influxdata/telegraf/telegraf --config $f &
done;

tail -f /var/log/telegraf-config-gen.log