package filter

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/processors"
)

var sampleConfig = `
# required_tags = ["tag1", "tag2"]
# required_fields = ["field1", "field2"]
`

type Filter struct {
	RequiredTags   []string `toml:"required_tags"`
	RequiredFields []string `toml:"required_fields"`
}

func New() *Filter {
	return &Filter{}
}

func (p *Filter) SampleConfig() string {
	return sampleConfig
}

func (p *Filter) Description() string {
	return ""
}

func (p *Filter) Apply(in ...telegraf.Metric) []telegraf.Metric {
	result := make([]telegraf.Metric, 0)
	for _, mt := range in {
		if p.hasRequired(mt) {
			result = append(result, mt)
		}
	}
	return result
}

func (p *Filter) hasRequired(metric telegraf.Metric) bool {
	for _, requiredTag := range p.RequiredTags {
		if !metric.HasTag(requiredTag) {
			return false
		}
	}
	for _, requiredField := range p.RequiredFields {
		if !metric.HasField(requiredField) {
			return false
		}
	}
	return true
}

func init() {
	processors.Add("filter", func() telegraf.Processor {
		return New()
	})
}
