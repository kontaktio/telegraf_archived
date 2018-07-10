package lastcalc

import (
	"log"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"github.com/influxdata/telegraf/plugins/processors"
)

var sampleConfig = `
# field_name = "name of field to process"
# out_field_name = "name of output field"
# tag_key = "tag to identify object"
# threshold = 5 # threshold in seconds (difference between next and prev value which triggers 1 on output)
`

type LastCalc struct {
	FieldName    string `toml:"field_name"`
	OutFieldName string `toml:"out_field_name"`
	TagKey       string `toml:"tag_key"`
	Threshold    int64  `toml:"threshold"`

	cache map[string]lastScanInfo
}

type lastScanInfo struct {
	lastValue          float64
	lastScanTime       time.Time
	lastZeroTime       time.Time
	lastZeroTimeExists bool
}

func New() *LastCalc {
	lastcalc := LastCalc{}

	lastcalc.Reset()
	return &lastcalc
}

func (p *LastCalc) Reset() {
	p.cache = make(map[string]lastScanInfo)
}

func (p *LastCalc) SampleConfig() string {
	return sampleConfig
}

func (p *LastCalc) Description() string {
	return ""
}

func (p *LastCalc) Apply(in ...telegraf.Metric) []telegraf.Metric {
	result := make([]telegraf.Metric, len(in))
	for idx, mt := range in {
		if mt.HasTag(p.TagKey) && mt.HasField(p.FieldName) {
			tag, _ := mt.GetTag(p.TagKey)
			fieldVal, _ := mt.GetField(p.FieldName)
			floatFieldVal, typeOk := fieldVal.(float64)
			if !typeOk {
				log.Printf("E! [processors.lastcalc] Invalid type of field %s", p.FieldName)
				return in //Wrong type
			}

			scanInfo, exists := p.cache[tag]
			if !exists {
				p.cache[tag] = lastScanInfo{floatFieldVal, mt.Time(), mt.Time(), false}
				result[idx] = p.copyAndReplaceField(mt, 0)
				continue
			}

			var newFieldValue int32

			if scanInfo.lastValue-float64(p.Threshold) > floatFieldVal &&
				mt.Time().Sub(scanInfo.lastZeroTime).Seconds() > float64(p.Threshold) {
				log.Printf("D! [processors.lastcalc] TrackingId: %s, Last val: %d, New val: %d\n", tag, int32(scanInfo.lastValue), int32(floatFieldVal))
				newFieldValue = 1
				p.cache[tag] = lastScanInfo{
					floatFieldVal,
					mt.Time(),
					mt.Time(),
					true,
				}
			} else {
				p.cache[tag] = lastScanInfo{
					floatFieldVal,
					mt.Time(),
					scanInfo.lastZeroTime,
					scanInfo.lastZeroTimeExists,
				}
				newFieldValue = 0
			}
			result[idx] = p.copyAndReplaceField(mt, newFieldValue)
		} else {
			result[idx] = mt
		}
	}
	return result
}

func (p *LastCalc) copyAndReplaceField(mt telegraf.Metric, newValue int32) telegraf.Metric {
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
	processors.Add("lastcalc", func() telegraf.Processor {
		return New()
	})
}
