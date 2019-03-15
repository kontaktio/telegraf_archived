package all

import (
	_ "github.com/influxdata/telegraf/plugins/inputs/http"
	_ "github.com/influxdata/telegraf/plugins/inputs/http_response"
	_ "github.com/influxdata/telegraf/plugins/inputs/httpjson"
	_ "github.com/influxdata/telegraf/plugins/inputs/influxdb"
	_ "github.com/influxdata/telegraf/plugins/inputs/influxdb_listener"
	_ "github.com/influxdata/telegraf/plugins/inputs/internal"
	_ "github.com/influxdata/telegraf/plugins/inputs/mqtt_consumer"
	_ "github.com/influxdata/telegraf/plugins/inputs/socket_listener"
)
