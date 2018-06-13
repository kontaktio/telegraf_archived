package all

import (
	_ "github.com/influxdata/telegraf/plugins/outputs/cloudwatch"
	_ "github.com/influxdata/telegraf/plugins/outputs/discard"
	_ "github.com/influxdata/telegraf/plugins/outputs/file"
	_ "github.com/influxdata/telegraf/plugins/outputs/http"
	_ "github.com/influxdata/telegraf/plugins/outputs/influxdb"
)
