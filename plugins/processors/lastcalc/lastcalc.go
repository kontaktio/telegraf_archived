package lastcalc

import (
	"log"
	"sync"
	"time"

	"github.com/influxdata/telegraf"
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
	cache        sync.Map
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
}

func (p *LastCalc) SampleConfig() string {
	return sampleConfig
}

func (p *LastCalc) Description() string {
	return ""
}

func (p *LastCalc) Apply(in ...telegraf.Metric) []telegraf.Metric {

	for _, mt := range in {
		if !mt.HasField(p.OutFieldName) && mt.HasField(p.FieldName) {
			tag, _ := mt.GetTag(p.TagKey)
			fieldVal, _ := mt.GetField(p.FieldName)
			floatFieldVal, typeOk := fieldVal.(float64)
			if !typeOk {
				log.Printf("E! [processors.lastcalc] Invalid type of field %s", p.FieldName)
				continue //Wrong type
			}

			item, exists := p.cache.Load(tag)
			if !exists {
				p.cache.Store(tag, lastScanInfo{floatFieldVal, mt.Time(), mt.Time(), false})
				mt.AddField(p.OutFieldName, 0)
				continue
			}

			scanInfo := item.(lastScanInfo)
			var newFieldValue int32

			if scanInfo.lastValue-float64(p.Threshold) > floatFieldVal &&
				mt.Time().Sub(scanInfo.lastZeroTime).Seconds() > float64(p.Threshold) {
				log.Printf("D! [processors.lastcalc] TrackingId: %s, Last val: %d, New val: %d\n", tag, int32(scanInfo.lastValue), int32(floatFieldVal))
				newFieldValue = 1
				p.cache.Store(tag, lastScanInfo{
					floatFieldVal,
					mt.Time(),
					mt.Time(),
					true,
				})
			} else {
				p.cache.Store(tag, lastScanInfo{
					floatFieldVal,
					mt.Time(),
					scanInfo.lastZeroTime,
					scanInfo.lastZeroTimeExists,
				})
				newFieldValue = 0
			}
			mt.AddField(p.OutFieldName, newFieldValue)
		}
	}

	return in
}

func init() {
	processors.Add("lastcalc", func() telegraf.Processor {
		return New()
	})
}
