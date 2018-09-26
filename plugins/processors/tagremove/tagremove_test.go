package tagremove

import (
	"testing"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"github.com/stretchr/testify/require"
)

const TagToRemove = "tagToRemove"
const TagToStay = "tagToStay"
const MetricName = "telemetry"

var metric1, _ = metric.New(
	MetricName,
	map[string]string{
		TagToRemove: "anything",
		TagToStay:   "efgh",
	},
	map[string]interface{}{
		"field": 3,
	},
	time.Now(),
)

func TestReplace(t *testing.T) {
	var tagremove TagRemove
	tagremove = *New()
	tagremove.Remove = []string{TagToRemove}

	var result = tagremove.Apply([]telegraf.Metric{metric1}...)
	require.Equal(t, 1, len(result))

	require.True(t, result[0].HasTag(TagToStay))
	require.False(t, result[0].HasTag(TagToRemove))
}
