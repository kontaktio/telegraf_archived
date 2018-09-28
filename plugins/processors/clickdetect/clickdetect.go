package clickdetect

import (
	"log"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
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

	cache map[string]float64
}

func New() *ClickDetect {
	lastcalc := ClickDetect{}

	lastcalc.Reset()
	return &lastcalc
}

func (p *ClickDetect) Reset() {
	p.cache = make(map[string]float64)
}

func (p *ClickDetect) SampleConfig() string {
	return sampleConfig
}

func (p *ClickDetect) Description() string {
	return ""
}

func (p *ClickDetect) Apply(in ...telegraf.Metric) []telegraf.Metric {
	result := make([]telegraf.Metric, len(in))
	for idx, mt := range in {
		if mt.HasTag(p.TagKey) && mt.HasField(p.FieldName) {
			tag, _ := mt.GetTag(p.TagKey)
			fieldVal, _ := mt.GetField(p.FieldName)
			floatField, typeOk := fieldVal.(float64)

			if !typeOk {
				log.Printf("E! [processors.clickdetect] Invalid type of field %s", p.FieldName)
				return in //Wrong type
			}

			lastValue, exists := p.cache[tag]

			p.cache[tag] = floatField
			if !exists || (lastValue == floatField) {
				result[idx] = p.copyAndReplaceField(mt, 0)
				continue
			}

			if floatField > lastValue {
				result[idx] = p.copyAndReplaceField(mt, int32(floatField-lastValue))
				continue
			} else if floatField < 10 && lastValue > 250 {
				diff := 256 - lastValue + floatField
				result[idx] = p.copyAndReplaceField(mt, int32(diff))
				continue
			}
			result[idx] = p.copyAndReplaceField(mt, 0)
		} else {
			result[idx] = mt
		}
	}
	return result
}

func (p *ClickDetect) copyAndReplaceField(mt telegraf.Metric, newValue int32) telegraf.Metric {
	newMetric, _ := metric.New(
		mt.Name(),
		mt.Tags(),
		make(map[string]interface{}),
		mt.Time())

	for k, v := range mt.Fields() {
		if k != p.FieldName {
			newMetric.AddField(k, v)
		} else {
			newMetric.AddField(p.OutFieldName, newValue)
		}
	}

	return newMetric
}

func init() {
	processors.Add("clickdetect", func() telegraf.Processor {
		return New()
	})
}
