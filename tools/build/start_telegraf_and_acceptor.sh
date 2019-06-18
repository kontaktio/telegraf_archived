#!/bin/sh

mkdir /config_generator
aws s3 cp s3://kontakt-telegraf-config/build-$1/telegraf.eventprocessor.$1.conf /

/usr/bin/telegraf -config /telegraf.eventprocessor.$1.conf