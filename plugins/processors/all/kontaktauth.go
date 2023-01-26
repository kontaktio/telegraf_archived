//go:build !custom || processors || processors.kontaktauth

package all

import _ "github.com/influxdata/telegraf/plugins/processors/kontaktauth" // register plugin
