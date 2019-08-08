package clickdetect

import (
	"log"
	"sync"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/processors"
)

var sampleConfig = `
# field_name = "name of field to process"
# out_field_name = "name of output field"
# tag_key = "tag to identify object"
`

type ClickDetect struct {
	FieldName    string `toml:"field_name"`
	OutFieldName string `toml:"out_field_name"`
	TagKey       string `toml:"tag_key"`

	cache sync.Map
}

func New() *ClickDetect {
	lastcalc := ClickDetect{}

	lastcalc.Reset()
	return &lastcalc
}

func (p *ClickDetect) Reset() {
}

func (p *ClickDetect) SampleConfig() string {
	return sampleConfig
}

func (p *ClickDetect) Description() string {
	return ""
}

func (p *ClickDetect) Apply(in ...telegraf.Metric) []telegraf.Metric {

	for _, mt := range in {
		if !mt.HasField(p.OutFieldName) && mt.HasField(p.FieldName) {
			tag, _ := mt.GetTag(p.TagKey)
			fieldVal, _ := mt.GetField(p.FieldName)
			currentValue, typeOk := fieldVal.(float64)

			if !typeOk {
				log.Printf("E! [processors.clickdetect] Invalid type of field %s", p.FieldName)
				continue //Wrong type
			}

			item, exists := p.cache.Load(tag)

			if !exists || (item.(float64) == currentValue) {
				mt.AddField(p.OutFieldName, 0)
				p.cache.Store(tag, currentValue)
				continue
			}

			lastValue := item.(float64)

			if currentValue > lastValue {
				mt.AddField(p.OutFieldName, int32(currentValue-lastValue))
				p.cache.Store(tag, currentValue)
				continue
			} else if currentValue < 10 && lastValue > 250 {
				diff := 256 - lastValue + currentValue
				mt.AddField(p.OutFieldName, int32(diff))
				p.cache.Store(tag, currentValue)
				continue
			}

			mt.AddField(p.OutFieldName, 0)
		}
	}
	return in
}

func init() {
	processors.Add("clickdetect", func() telegraf.Processor {
		return New()
	})
}
