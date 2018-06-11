package lastcalc

import (
	"testing"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"github.com/stretchr/testify/require"
)

const FieldToChange = "fieldToChange"
const ChangedField = "changedField"
const MetricName = "telemetry"

var metric1, _ = metric.New(
	MetricName,
	map[string]string{"topic": "abcd"},
	map[string]interface{}{
		FieldToChange: float64(45),
	},
	time.Now().Add(time.Duration(-2)),
)

var metric2, _ = metric.New(
	MetricName,
	map[string]string{"topic": "abcd"},
	map[string]interface{}{
		FieldToChange: float64(2),
	},
	time.Now(),
)

var metric3, _ = metric.New(
	MetricName,
	map[string]string{"topic": "abcd"},
	map[string]interface{}{
		FieldToChange: float64(2),
	},
	time.Now().Add(time.Duration(2)),
)

func TestCreate(t *testing.T) {
	var lastc LastCalc
	lastc = *New()
	lastc.FieldName = FieldToChange
	lastc.OutFieldName = ChangedField
	lastc.Threshold = 1
	lastc.TagKey = "topic"
	lastc.Reset()

	var result = lastc.Apply([]telegraf.Metric{metric1}...)
	require.Equal(t, 1, len(result))
	var changedFieldValue, changedFieldExists = result[0].GetField(ChangedField)
	require.True(t, changedFieldExists)
	require.Equal(t, int64(0), changedFieldValue)
	require.False(t, result[0].HasField(FieldToChange))

	result = lastc.Apply([]telegraf.Metric{metric2}...)
	require.Equal(t, 1, len(result))
	changedFieldValue, changedFieldExists = result[0].GetField(ChangedField)
	require.True(t, changedFieldExists)
	require.Equal(t, int64(1), changedFieldValue)
	require.False(t, result[0].HasField(FieldToChange))

	result = lastc.Apply([]telegraf.Metric{metric3}...)
	require.Equal(t, 1, len(result))
	changedFieldValue, changedFieldExists = result[0].GetField(ChangedField)
	require.True(t, changedFieldExists)
	require.Equal(t, int64(0), changedFieldValue)
	require.False(t, result[0].HasField(FieldToChange))
}
