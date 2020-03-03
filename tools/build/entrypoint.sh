#!/bin/sh
set -e

gomplate -f /telegraf.conf.tpl -o /telegraf.conf
/usr/bin/telegraf -config /telegraf.conf
