package tagrename

import (
	"testing"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"github.com/stretchr/testify/require"
)

const TagToChange = "tagToChange"
const TagToChangeValue = "tagToChangeValue"
const TagToStay = "tagToStay"
const ChangedTag = "changedTag"
const MetricName = "telemetry"

var metric1, _ = metric.New(
	MetricName,
	map[string]string{
		TagToChange: TagToChangeValue,
		TagToStay:   "efgh",
	},
	map[string]interface{}{
		"field": 3,
	},
	time.Now(),
)

func TestReplace(t *testing.T) {
	var tagrename TagRename
	tagrename = *New()
	tagrename.Renames = map[string]string{TagToChange: ChangedTag}

	var result = tagrename.Apply([]telegraf.Metric{metric1}...)
	require.Equal(t, 1, len(result))

	require.True(t, result[0].HasTag(TagToStay))
	require.False(t, result[0].HasTag(TagToChange))
	tagValue, tagExists := result[0].GetTag(ChangedTag)
	require.True(t, tagExists)
	require.Equal(t, TagToChangeValue, tagValue)
}
