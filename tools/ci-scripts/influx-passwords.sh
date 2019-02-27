#!/bin/bash

root_password=
new_password=
influx_host=testinflux.kontakt.io

while [[ "$1" != "" ]]; do
    case $1 in
        -p | --root-password )  shift
                                root_password=$1
                                ;;
        -n | --new-password )   shift
		                        new_password=$1
                                ;;
        -h | --influx-host )	shift
		                        influx_host=$1
                                ;;
        * )                     echo "(-p | --root_password) (-n | --new_password) [-h | --influx_host]"
                                exit 1
    esac
    shift
done

influx -host ${influx_host} -username root -password "$root_password" -execute 'SHOW USERS' -format 'csv' \
 | grep 'false$' | cut -d ',' -f1 | sed -e 's/^/SET PASSWORD FOR "/' | sed -e 's/$/"/' | sed -e "s/$/ = '$new_password'/" \
 > reset-passwords.influx

influx -host ${influx_host} -username root -password "$root_password" -execute "`cat reset-passwords.influx`"
