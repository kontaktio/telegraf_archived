#!/bin/sh

mkdir /config_generator
aws s3 cp $configPath /configFile

/usr/bin/telegraf -config /configFile
