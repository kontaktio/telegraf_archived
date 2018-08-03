package trackingidresolver

import (
	"testing"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"github.com/stretchr/testify/require"
)

const TagName = "ble_proximityuuid"
const NewTagName = "trackingId"
const OtherTagName = "RandomField"
const MetricName = "telemetry"

var mt, _ = metric.New(
	MetricName,
	map[string]string{
		TagName: "426f4e4f-0000-0000-0000-000000000000",
	},
	map[string]interface{}{
		"topic": "abcd",
	},
	time.Now(),
)
var noFieldMt, _ = metric.New(
	MetricName,
	map[string]string{
		OtherTagName: "426f4e4f-0000-0000-0000-000000000000",
	},
	map[string]interface{}{
		"topic": "abcd",
	},
	time.Now(),
)

func TestResolveProximity(t *testing.T) {
	trackingidresolver := New()
	trackingidresolver.TagName = TagName
	trackingidresolver.NewTagName = NewTagName

	var result = trackingidresolver.Apply([]telegraf.Metric{mt}...)
	require.Equal(t, 1, len(result))

	var changedTagVal, changedTagExists = result[0].GetTag(NewTagName)
	require.True(t, changedTagExists)
	require.Equal(t, "BoNO", changedTagVal)
	require.False(t, result[0].HasField(TagName))
}

func TestDoNothingWhenNoField(t *testing.T) {
	trackingidresolver := New()
	trackingidresolver.TagName = TagName
	trackingidresolver.NewTagName = NewTagName

	var result = trackingidresolver.Apply([]telegraf.Metric{noFieldMt}...)
	require.Equal(t, 1, len(result))

	require.False(t, result[0].HasTag(TagName))
	require.False(t, result[0].HasTag(NewTagName))
	require.True(t, result[0].HasTag(OtherTagName))
}
