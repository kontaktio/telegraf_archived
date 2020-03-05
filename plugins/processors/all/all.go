package all

import (
	_ "github.com/influxdata/telegraf/plugins/processors/clickdetect"
	_ "github.com/influxdata/telegraf/plugins/processors/converter"
	_ "github.com/influxdata/telegraf/plugins/processors/filter"
	_ "github.com/influxdata/telegraf/plugins/processors/kioparser"
	_ "github.com/influxdata/telegraf/plugins/processors/kontaktauth"
	_ "github.com/influxdata/telegraf/plugins/processors/lastcalc"
	_ "github.com/influxdata/telegraf/plugins/processors/override"
	_ "github.com/influxdata/telegraf/plugins/processors/parser"
	_ "github.com/influxdata/telegraf/plugins/processors/printer"
	_ "github.com/influxdata/telegraf/plugins/processors/regex"
	_ "github.com/influxdata/telegraf/plugins/processors/rename"
	_ "github.com/influxdata/telegraf/plugins/processors/strings"
	_ "github.com/influxdata/telegraf/plugins/processors/topk"
	_ "github.com/influxdata/telegraf/plugins/processors/trackingidresolver"
)
