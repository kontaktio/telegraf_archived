package fieldremove

import (
	"testing"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"github.com/stretchr/testify/require"
)

const FieldToRemove = "toremove"
const FieldToStay = "tostay"
const MetricName = "telemetry"

var metric1, _ = metric.New(
	MetricName,
	map[string]string{
		"Tag": "anything",
	},
	map[string]interface{}{
		FieldToRemove: 3,
		FieldToStay:   "abcd",
	},
	time.Now(),
)

func TestReplace(t *testing.T) {
	var fieldremove FieldRemove
	fieldremove = *New()
	fieldremove.Remove = []string{FieldToRemove}

	var result = fieldremove.Apply([]telegraf.Metric{metric1}...)
	require.Equal(t, 1, len(result))

	require.True(t, result[0].HasField(FieldToStay))
	require.False(t, result[0].HasField(FieldToRemove))
}
