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

func TestNotEmit1WhenTimeDiffTooSmall(t *testing.T) {
	lastc := *New()
	lastc.FieldName = FieldToChange
	lastc.OutFieldName = ChangedField
	lastc.Threshold = 1
	lastc.TagKey = "topic"
	lastc.Reset()

	var metric1Time, _ = time.ParseDuration("0m")
	var metric2Time, _ = time.ParseDuration("500ms")

	var metric1, _ = metric.New(
		MetricName,
		map[string]string{"topic": "abcd"},
		map[string]interface{}{
			FieldToChange: float64(45),
		},
		time.Now().Add(metric1Time),
	)

	var metric2, _ = metric.New(
		MetricName,
		map[string]string{"topic": "abcd"},
		map[string]interface{}{
			FieldToChange: float64(2),
		},
		time.Now().Add(metric2Time),
	)

	var result = lastc.Apply([]telegraf.Metric{metric1}...)
	require.Equal(t, 1, len(result))
	var changedFieldValue, changedFieldExists = result[0].GetField(ChangedField)
	require.True(t, changedFieldExists)
	require.Equal(t, int64(0), changedFieldValue)

	result = lastc.Apply([]telegraf.Metric{metric2}...)
	require.Equal(t, 1, len(result))
	changedFieldValue, changedFieldExists = result[0].GetField(ChangedField)
	require.True(t, changedFieldExists)
	require.Equal(t, int64(0), changedFieldValue)
}

func TestCreate(t *testing.T) {
	var lastc LastCalc
	lastc = *New()
	lastc.FieldName = FieldToChange
	lastc.OutFieldName = ChangedField
	lastc.Threshold = 1
	lastc.TagKey = "topic"
	lastc.Reset()

	var metric1Time, _ = time.ParseDuration("-1m")
	var metric2Time, _ = time.ParseDuration("0m")
	var metric3Time, _ = time.ParseDuration("1m")

	var metric1, _ = metric.New(
		MetricName,
		map[string]string{"topic": "abcd"},
		map[string]interface{}{
			FieldToChange: float64(45),
		},
		time.Now().Add(metric1Time),
	)

	var metric2, _ = metric.New(
		MetricName,
		map[string]string{"topic": "abcd"},
		map[string]interface{}{
			FieldToChange: float64(2),
		},
		time.Now().Add(metric2Time),
	)

	var metric3, _ = metric.New(
		MetricName,
		map[string]string{"topic": "abcd"},
		map[string]interface{}{
			FieldToChange: float64(2),
		},
		time.Now().Add(metric3Time),
	)

	var result = lastc.Apply([]telegraf.Metric{metric1}...)
	require.Equal(t, 1, len(result))
	var changedFieldValue, changedFieldExists = result[0].GetField(ChangedField)
	require.True(t, changedFieldExists)
	require.Equal(t, int64(0), changedFieldValue)

	result = lastc.Apply([]telegraf.Metric{metric2}...)
	require.Equal(t, 1, len(result))
	changedFieldValue, changedFieldExists = result[0].GetField(ChangedField)
	require.True(t, changedFieldExists)
	require.Equal(t, int64(1), changedFieldValue)

	result = lastc.Apply([]telegraf.Metric{metric3}...)
	require.Equal(t, 1, len(result))
	changedFieldValue, changedFieldExists = result[0].GetField(ChangedField)
	require.True(t, changedFieldExists)
	require.Equal(t, int64(0), changedFieldValue)
}
