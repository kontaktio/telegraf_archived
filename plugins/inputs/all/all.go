package all

import (
	_ "github.com/influxdata/telegraf/plugins/inputs/cloudwatch"
	_ "github.com/influxdata/telegraf/plugins/inputs/exec"
	_ "github.com/influxdata/telegraf/plugins/inputs/http"
	_ "github.com/influxdata/telegraf/plugins/inputs/http_listener"
	_ "github.com/influxdata/telegraf/plugins/inputs/http_response"
	_ "github.com/influxdata/telegraf/plugins/inputs/httpjson"
	_ "github.com/influxdata/telegraf/plugins/inputs/influxdb"
	_ "github.com/influxdata/telegraf/plugins/inputs/internal"
	_ "github.com/influxdata/telegraf/plugins/inputs/kapacitor"
	_ "github.com/influxdata/telegraf/plugins/inputs/mqtt_consumer"
	_ "github.com/influxdata/telegraf/plugins/inputs/net_response"
	_ "github.com/influxdata/telegraf/plugins/inputs/postgresql"
	_ "github.com/influxdata/telegraf/plugins/inputs/postgresql_extensible"
	_ "github.com/influxdata/telegraf/plugins/inputs/rabbitmq"
	_ "github.com/influxdata/telegraf/plugins/inputs/redis"
	_ "github.com/influxdata/telegraf/plugins/inputs/tomcat"
	_ "github.com/influxdata/telegraf/plugins/inputs/internal"
)
