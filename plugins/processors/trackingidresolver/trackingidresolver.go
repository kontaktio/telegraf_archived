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
	for _, mt := range in {
		tagVal, _ := mt.GetTag(p.TagName)
		proximity, err := uuid.FromString(tagVal)
		if err != nil {
			log.Printf("E! [processors.trackingidresolver] Proximity is of invalid format: %s", tagVal)
			continue
		}
		proximityBytes := proximity.Bytes()
		len := bytes.IndexByte(proximityBytes, 0)
		trackingID := string(proximityBytes[:len])

		mt.AddTag(p.NewTagName, trackingID)
	}
	return in
}

func init() {
	processors.Add("trackingidresolver", func() telegraf.Processor {
		return New()
	})
}
