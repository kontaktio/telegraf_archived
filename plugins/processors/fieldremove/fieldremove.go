package fieldremove

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/processors"
)

var sampleConfig = `
# [processors.fieldremove]
# remove=['a','b','c']
`

type FieldRemove struct {
	Remove []string
}

func New() *FieldRemove {
	return &FieldRemove{}
}

func (t *FieldRemove) SampleConfig() string {
	return sampleConfig
}

func (t *FieldRemove) Description() string {
	return ""
}

func (t *FieldRemove) Apply(in ...telegraf.Metric) []telegraf.Metric {
	for _, mt := range in {
		for _, tag := range t.Remove {
			if mt.HasField(tag) {
				mt.RemoveField(tag)
			}
		}
	}
	return in
}

func init() {
	processors.Add("fieldremove", func() telegraf.Processor {
		return New()
	})
}
