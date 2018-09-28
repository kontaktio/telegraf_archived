package clickdetect

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

func TestEmit0WhenClickIdEqual(t *testing.T) {
	cd := *New()
	cd.FieldName = FieldToChange
	cd.OutFieldName = ChangedField
	cd.TagKey = "topic"
	cd.Reset()

	var metric2Time, _ = time.ParseDuration("500ms")

	var metric1, _ = metric.New(
		MetricName,
		map[string]string{
			"topic":      "abcd",
			"anotherTag": "q",
		},
		map[string]interface{}{
			FieldToChange:  float64(45),
			"anotherField": float64(55),
		},
		time.Now(),
	)

	var metric2, _ = metric.New(
		MetricName,
		map[string]string{
			"topic":      "abcd",
			"anotherTag": "q",
		},
		map[string]interface{}{
			FieldToChange:  float64(45),
			"anotherField": float64(55),
		},
		time.Now().Add(metric2Time),
	)

	var result = cd.Apply([]telegraf.Metric{metric1}...)
	var changedFieldValue, changedFieldExists = result[0].GetField(ChangedField)
	require.True(t, changedFieldExists)
	require.Equal(t, int64(0), changedFieldValue)
	require.False(t, result[0].HasField(FieldToChange))

	result = cd.Apply([]telegraf.Metric{metric2}...)
	changedFieldValue, changedFieldExists = result[0].GetField(ChangedField)
	require.True(t, changedFieldExists)
	require.Equal(t, int64(0), changedFieldValue)
	require.False(t, result[0].HasField(FieldToChange))
}

func TestNotEmit0WhenClickIdGreater(t *testing.T) {
	cd := *New()
	cd.FieldName = FieldToChange
	cd.OutFieldName = ChangedField
	cd.TagKey = "topic"
	cd.Reset()

	var metric2Time, _ = time.ParseDuration("500ms")

	var metric1, _ = metric.New(
		MetricName,
		map[string]string{
			"topic":      "abcd",
			"anotherTag": "q",
		},
		map[string]interface{}{
			FieldToChange:  float64(45),
			"anotherField": float64(55),
		},
		time.Now(),
	)

	var metric2, _ = metric.New(
		MetricName,
		map[string]string{
			"topic":      "abcd",
			"anotherTag": "q",
		},
		map[string]interface{}{
			FieldToChange:  float64(51),
			"anotherField": float64(55),
		},
		time.Now().Add(metric2Time),
	)

	var result = cd.Apply([]telegraf.Metric{metric1}...)
	var changedFieldValue, changedFieldExists = result[0].GetField(ChangedField)
	require.True(t, changedFieldExists)
	require.Equal(t, int64(0), changedFieldValue)
	require.False(t, result[0].HasField(FieldToChange))

	result = cd.Apply([]telegraf.Metric{metric2}...)
	changedFieldValue, changedFieldExists = result[0].GetField(ChangedField)
	require.True(t, changedFieldExists)
	require.Equal(t, int64(6), changedFieldValue)
	require.False(t, result[0].HasField(FieldToChange))
}

func TestNotEmit0WhenClickIdRolledBack(t *testing.T) {
	cd := *New()
	cd.FieldName = FieldToChange
	cd.OutFieldName = ChangedField
	cd.TagKey = "topic"
	cd.Reset()

	var metric2Time, _ = time.ParseDuration("500ms")

	var metric1, _ = metric.New(
		MetricName,
		map[string]string{
			"topic":      "abcd",
			"anotherTag": "q",
		},
		map[string]interface{}{
			FieldToChange:  float64(252),
			"anotherField": float64(55),
		},
		time.Now(),
	)

	var metric2, _ = metric.New(
		MetricName,
		map[string]string{
			"topic":      "abcd",
			"anotherTag": "q",
		},
		map[string]interface{}{
			FieldToChange:  float64(4),
			"anotherField": float64(55),
		},
		time.Now().Add(metric2Time),
	)

	var result = cd.Apply([]telegraf.Metric{metric1}...)
	var changedFieldValue, changedFieldExists = result[0].GetField(ChangedField)
	require.True(t, changedFieldExists)
	require.Equal(t, int64(0), changedFieldValue)
	require.False(t, result[0].HasField(FieldToChange))

	result = cd.Apply([]telegraf.Metric{metric2}...)
	changedFieldValue, changedFieldExists = result[0].GetField(ChangedField)
	require.True(t, changedFieldExists)
	require.Equal(t, int64(8), changedFieldValue)
	require.False(t, result[0].HasField(FieldToChange))
}

func TestNotEmit0WhenClickIdLower(t *testing.T) {
	cd := *New()
	cd.FieldName = FieldToChange
	cd.OutFieldName = ChangedField
	cd.TagKey = "topic"
	cd.Reset()

	var metric2Time, _ = time.ParseDuration("500ms")

	var metric1, _ = metric.New(
		MetricName,
		map[string]string{
			"topic":      "abcd",
			"anotherTag": "q",
		},
		map[string]interface{}{
			FieldToChange:  float64(123),
			"anotherField": float64(55),
		},
		time.Now(),
	)

	var metric2, _ = metric.New(
		MetricName,
		map[string]string{
			"topic":      "abcd",
			"anotherTag": "q",
		},
		map[string]interface{}{
			FieldToChange:  float64(110),
			"anotherField": float64(55),
		},
		time.Now().Add(metric2Time),
	)

	var result = cd.Apply([]telegraf.Metric{metric1}...)
	var changedFieldValue, changedFieldExists = result[0].GetField(ChangedField)
	require.True(t, changedFieldExists)
	require.Equal(t, int64(0), changedFieldValue)
	require.False(t, result[0].HasField(FieldToChange))

	result = cd.Apply([]telegraf.Metric{metric2}...)
	changedFieldValue, changedFieldExists = result[0].GetField(ChangedField)
	require.True(t, changedFieldExists)
	require.Equal(t, int64(0), changedFieldValue)
	require.False(t, result[0].HasField(FieldToChange))
}
