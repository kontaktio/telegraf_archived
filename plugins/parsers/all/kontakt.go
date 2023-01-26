//go:build !custom || parsers || parsers.kontakt

package all

import (
	_ "github.com/influxdata/telegraf/plugins/parsers/kontakt" // register plugin
)
