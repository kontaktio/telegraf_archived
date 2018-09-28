package all

import (
	_ "github.com/influxdata/telegraf/plugins/processors/clickdetect"
	_ "github.com/influxdata/telegraf/plugins/processors/converter"
	_ "github.com/influxdata/telegraf/plugins/processors/fieldremove"
	_ "github.com/influxdata/telegraf/plugins/processors/lastcalc"
	_ "github.com/influxdata/telegraf/plugins/processors/override"
	_ "github.com/influxdata/telegraf/plugins/processors/printer"
	_ "github.com/influxdata/telegraf/plugins/processors/regex"
	_ "github.com/influxdata/telegraf/plugins/processors/tagremove"
	_ "github.com/influxdata/telegraf/plugins/processors/tagrename"
	_ "github.com/influxdata/telegraf/plugins/processors/topk"
	_ "github.com/influxdata/telegraf/plugins/processors/trackingidresolver"
)
