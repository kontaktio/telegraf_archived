package printer

import (
	"fmt"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/processors"
	"github.com/influxdata/telegraf/plugins/serializers"
	"github.com/influxdata/telegraf/plugins/serializers/influx"
)

type Printer struct {
	Tags       map[string][]string
	serializer serializers.Serializer
}

var sampleConfig = `
	## Tags with values for which metric should be printed 
	# [processors.printer.tags]
    #   tag1Name = ["tag1_value", "tag1_other_value"]
`

func (p *Printer) SampleConfig() string {
	return sampleConfig
}

func (p *Printer) Description() string {
	return "Print all metrics that pass through this filter."
}

func (p *Printer) Apply(in ...telegraf.Metric) []telegraf.Metric {
	for _, metric := range in {
		if p.Tags == nil || p.shouldPrint(metric) {
			octets, err := p.serializer.Serialize(metric)
			if err != nil {
				continue
			}
			fmt.Printf("%s", octets)
		}
	}
	return in
}

func (p *Printer) shouldPrint(metric telegraf.Metric) bool {
	for wantedName, wantedValues := range p.Tags {
		tagValue, ok := metric.GetTag(wantedName)
		if !ok {
			return false
		}
		found := false
		for _, wantedValue := range wantedValues {
			if wantedValue == tagValue {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func init() {
	processors.Add("printer", func() telegraf.Processor {
		return &Printer{
			serializer: influx.NewSerializer(),
		}
	})
}
