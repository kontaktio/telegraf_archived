package tagrename

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/processors"
)

var sampleConfig = `
# [processors.tagrename.renames]
# tag1 = "renamedtag1"
`

type TagRename struct {
	Renames map[string]string
}

func New() *TagRename {
	return &TagRename{}
}

func (t *TagRename) SampleConfig() string {
	return sampleConfig
}

func (t *TagRename) Description() string {
	return ""
}

func (t *TagRename) Apply(in ...telegraf.Metric) []telegraf.Metric {
	for _, mt := range in {
		for tag, replacement := range t.Renames {
			if mt.HasTag(tag) {
				tagValue, _ := mt.GetTag(tag)
				mt.AddTag(replacement, tagValue)
				mt.RemoveTag(tag)
			}
		}
	}
	return in
}

func init() {
	processors.Add("tagrename", func() telegraf.Processor {
		return New()
	})
}
