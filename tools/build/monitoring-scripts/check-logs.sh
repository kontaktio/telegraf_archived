#!/bin/sh
#180 = last 15 minutes, (15 * 60 seconds) / 5 seconds flush buffer which generate log every 5s
tail -n 180 /var/log/telegraf-config-gen.log | grep "wrote batch of" | wc -l | xargs /exit.sh && pm2 restart all || echo "ok"
