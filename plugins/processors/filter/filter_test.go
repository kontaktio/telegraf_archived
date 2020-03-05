package filter

import (
	"github.com/influxdata/telegraf/metric"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const TagName = "qwert"
const FieldName = "abcdef"
const MetricName = "telemetry"

func TestFilter(t *testing.T) {
	f := New()
	f.RequiredTags = []string{TagName}
	f.RequiredFields = []string{FieldName}

	var mt1, _ = metric.New(
		MetricName,
		map[string]string{
			TagName: "426f4e4f-0000-0000-0000-000000000000",
		},
		map[string]interface{}{
			FieldName: "abcd",
		},
		time.Now(),
	)

	assert.Equal(t, 1, len(f.Apply(mt1)))

	var mt2, _ = metric.New(
		MetricName,
		map[string]string{
			TagName: "426f4e4f-0000-0000-0000-000000000000",
		},
		map[string]interface{}{
			"OtherField": "abcd",
		},
		time.Now(),
	)
	assert.Equal(t, 0, len(f.Apply(mt2)))

	var mt3, _ = metric.New(
		MetricName,
		map[string]string{
			"otherTag": "426f4e4f-0000-0000-0000-000000000000",
		},
		map[string]interface{}{
			FieldName: "abcd",
		},
		time.Now(),
	)

	assert.Equal(t, 0, len(f.Apply(mt3)))
}
