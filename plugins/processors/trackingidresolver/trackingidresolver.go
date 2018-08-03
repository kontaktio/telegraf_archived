package trackingidresolver

import (
	"bytes"
	"log"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"github.com/influxdata/telegraf/plugins/processors"
	"github.com/satori/go.uuid"
)

var sampleConfig = `
# tag_name = "name of tag to process"
# new_tag_name = "name to save resolved trackingId"
`

type TrackingIdResolver struct {
	TagName    string `toml:"tag_name"`
	NewTagName string `toml:"new_tag_name"`
}

func New() *TrackingIdResolver {
	return &TrackingIdResolver{}
}

func (p *TrackingIdResolver) SampleConfig() string {
	return sampleConfig
}

func (p *TrackingIdResolver) Description() string {
	return ""
}

func (p *TrackingIdResolver) Apply(in ...telegraf.Metric) []telegraf.Metric {
	result := make([]telegraf.Metric, len(in))
	for idx, mt := range in {
		if mt.HasTag(p.TagName) {
			tagVal, _ := mt.GetTag(p.TagName)
			proximity, err := uuid.FromString(tagVal)
			if err != nil {
				log.Printf("E! [processors.trackingidresolver] Proximity is of invalid format: %s", tagVal)
				return in //Wrong type
			}
			proximityBytes := proximity.Bytes()
			len := bytes.IndexByte(proximityBytes, 0)
			trackingID := string(proximityBytes[:len])
			result[idx] = p.copyAndReplaceField(mt, trackingID)
		} else {
			log.Printf("W! [processors.trackingidresolver] No Tag with name %s", p.TagName)
			result[idx] = mt
		}
	}
	return result
}

func (p *TrackingIdResolver) copyAndReplaceField(mt telegraf.Metric, newValue string) telegraf.Metric {
	newMetric, _ := metric.New(
		mt.Name(),
		make(map[string]string),
		mt.Fields(),
		mt.Time())

	for k, v := range mt.Tags() {
		if k != p.TagName {
			newMetric.AddTag(k, v)
		} else {
			newMetric.AddTag(p.NewTagName, newValue)
		}
	}

	return newMetric
}

func init() {
	processors.Add("trackingidresolver", func() telegraf.Processor {
		return New()
	})
}
