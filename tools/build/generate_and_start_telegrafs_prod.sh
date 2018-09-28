#TODO Review PS->MK: może zróbmy troche porządek z plikami, podzielmy na katalogi o konkretnym przeznaczeniu
# - Docker - to co potrzebne do build
# - telegraf-run-cofig, telegraf-run-config-prod, telegraf-run-config-test - to co należy skopiować na S3, co jest czytane prze telegraf w czasie uruchamiania
# - ci-scripts - kopia poleceń z jenkins, przykładowo pliki deployment_script_*
# - kapacitor - to co jest do wgrania do kapacitora, chociaz nie wiem czy nie lepiej zrobić fork kapacitora i tam umieścić w podobnej konwencji...

#!/bin/sh
cd /etc/telegraf
python /config_generator/telegraf_stream_config_generate.py --api-key $API_KEY --influxdb-url ${INFLUXDB_URL:=http://influx.kontakt.io} --influxdb-port ${INFLUXDB_PORT:=8086} --influxdb-username $INFLUXDB_USERNAME --influxdb-password $INFLUXDB_PASSWORD --api-url ${API_URL:=http://api.kontakt.io} --data-collection-interval ${DATA_COLLECTION_INTERVAL:=5} --flush-interval ${FLUSH_INTERVAL:=5} --flush-jitter ${FLUSH_JITTER:=2} --debug ${DEBUG_ENABLED:=true} --log-file ${LOG_FILE_NAME:=/var/log/telegraf-config-gen.log} --mqtt-url ${MQTT_URL:=ssl://ovs.kontakt.io:8083} --telemetry-types ${TELEMETRY_TYPES:=health sensor accelerometer button} --streams-per-telegraf ${STREAMS_PER_TELEGRAF:=100} $(if [ ! -z $VENUE_ID ]; then echo "--api-venue-id $VENUE_ID"; else echo ""; fi;)
for f in /etc/telegraf/telegraf.stream.conf.*;
do
	/go/src/github.com/influxdata/telegraf/telegraf --config $f &
done;

python /config_generator/telegraf_location_config_generate.py --api-key $API_KEY --kapacitor-url ${KAPACITOR_URL:=http://influx.kontakt.io:9090} --kapacitor-user ${KAPACITOR_USER:=kontaktio} --kapacitor-pass ${KAPACITOR_PASS:=notthepassword} --influxdb-url ${INFLUXDB_URL:=http://influx.kontakt.io} --influxdb-port ${INFLUXDB_PORT:=8086} --influxdb-username $INFLUXDB_USERNAME --influxdb-password $INFLUXDB_PASSWORD --api-url ${API_URL:=http://api.kontakt.io} --data-collection-interval ${DATA_COLLECTION_INTERVAL:=5} --flush-interval ${FLUSH_INTERVAL:=5} --flush-jitter ${FLUSH_JITTER:=2} --debug ${DEBUG_ENABLED:=true} --log-file ${LOG_FILE_NAME:=/var/log/telegraf-config-gen.log} --tx-power ${TX_POWER:=-77} $(if [ ! -z $VENUE_ID ]; then echo "--api-venue-id $VENUE_ID"; else echo ""; fi;)
for f in /etc/telegraf/telegraf.location.conf.*;
do
	/go/src/github.com/influxdata/telegraf/telegraf --config $f &
done;

tail -f /var/log/telegraf-config-gen.log