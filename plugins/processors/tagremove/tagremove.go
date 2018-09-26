package tagremove

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/processors"
)

var sampleConfig = `
# [processors.tagremove]
# remove=['a','b','c']
`

type TagRemove struct {
	Remove []string
}

func New() *TagRemove {
	return &TagRemove{}
}

func (t *TagRemove) SampleConfig() string {
	return sampleConfig
}

func (t *TagRemove) Description() string {
	return ""
}

func (t *TagRemove) Apply(in ...telegraf.Metric) []telegraf.Metric {
	for _, mt := range in {
		for _, tag := range t.Remove {
			if mt.HasTag(tag) {
				mt.RemoveTag(tag)
			}
		}
	}
	return in
}

func init() {
	processors.Add("tagremove", func() telegraf.Processor {
		return New()
	})
}
