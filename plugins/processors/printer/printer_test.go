package printer

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"github.com/influxdata/telegraf/plugins/serializers"
	"github.com/influxdata/telegraf/plugins/serializers/influx"
	"testing"
	"time"
)

func TestPrintingNoTags(t *testing.T) {
	serializer := getTestSerializer()
	printer := Printer{
		Tags:       nil,
		serializer: &serializer,
	}
	builder := metric.NewBuilder()
	builder.SetTime(time.Now())
	m, _ := builder.Metric()
	printer.Apply(m)
	if len(serializer.saved) != 1 {
		t.Fatalf("number of serialized, expected 1, got %v", len(serializer.saved))
	}
}

func TestPrintingHasTagWithValue(t *testing.T) {
	serializer := getTestSerializer()
	printer := Printer{
		Tags: map[string][]string{
			"tag1": {"value1"},
		},
		serializer: &serializer,
	}
	builder := metric.NewBuilder()
	builder.SetTime(time.Now())
	builder.AddTag("tag1", "value1")
	m, _ := builder.Metric()
	printer.Apply(m)
	if len(serializer.saved) != 1 {
		t.Fatalf("number of serialized, expected 1, got %v", len(serializer.saved))
	}
}

func TestPrintingHasMultipleTagsWithValue(t *testing.T) {
	serializer := getTestSerializer()
	printer := Printer{
		Tags: map[string][]string{
			"tag1": {"value1"},
			"tag2": {"value2"},
		},
		serializer: &serializer,
	}
	builder := metric.NewBuilder()
	builder.SetTime(time.Now())
	builder.AddTag("tag1", "value1")
	builder.AddTag("tag2", "value2")
	m, _ := builder.Metric()
	printer.Apply(m)
	if len(serializer.saved) != 1 {
		t.Fatalf("number of serialized, expected 1, got %v", len(serializer.saved))
	}
}

func TestPrintingHasTagWithDifferentValue(t *testing.T) {
	serializer := getTestSerializer()
	printer := Printer{
		Tags: map[string][]string{
			"tag1": {"value1"},
		},
		serializer: &serializer,
	}
	builder := metric.NewBuilder()
	builder.SetTime(time.Now())
	builder.AddTag("tag1", "value2")
	m, _ := builder.Metric()
	printer.Apply(m)
	if len(serializer.saved) != 0 {
		t.Fatalf("number of serialized, expected 0, got %v", len(serializer.saved))
	}
}

func TestPrintingHasTagWithValueNotInCollection(t *testing.T) {
	serializer := getTestSerializer()
	printer := Printer{
		Tags: map[string][]string{
			"tag1": {"value1", "value2"},
		},
		serializer: &serializer,
	}
	builder := metric.NewBuilder()
	builder.SetTime(time.Now())
	builder.AddTag("tag1", "value3")
	m, _ := builder.Metric()
	printer.Apply(m)
	if len(serializer.saved) != 0 {
		t.Fatalf("number of serialized, expected 0, got %v", len(serializer.saved))
	}
}

func TestPrintingHasTagOnlyOneOfTags(t *testing.T) {
	serializer := getTestSerializer()
	printer := Printer{
		Tags: map[string][]string{
			"tag1": {"value11", "value12"},
			"tag2": {"value21", "value22"},
		},
		serializer: &serializer,
	}
	builder := metric.NewBuilder()
	builder.SetTime(time.Now())
	builder.AddTag("tag1", "value11")
	m, _ := builder.Metric()
	printer.Apply(m)
	if len(serializer.saved) != 0 {
		t.Fatalf("number of serialized, expected 0, got %v", len(serializer.saved))
	}
}

func TestPrintingHasTagOnlyOneOfTagValuesCorrect(t *testing.T) {
	serializer := getTestSerializer()
	printer := Printer{
		Tags: map[string][]string{
			"tag1": {"value11", "value12"},
			"tag2": {"value21", "value22"},
		},
		serializer: &serializer,
	}
	builder := metric.NewBuilder()
	builder.SetTime(time.Now())
	builder.AddTag("tag1", "value11")
	builder.AddTag("tag2", "value23")
	m, _ := builder.Metric()
	printer.Apply(m)
	if len(serializer.saved) != 0 {
		t.Fatalf("number of serialized, expected 0, got %v", len(serializer.saved))
	}
}

type testSerializer struct {
	saved           []telegraf.Metric
	innerSerializer serializers.Serializer
}

func getTestSerializer() testSerializer {
	return testSerializer{
		saved:           make([]telegraf.Metric, 0),
		innerSerializer: influx.NewSerializer(),
	}
}

func (t *testSerializer) Serialize(metric telegraf.Metric) ([]byte, error) {
	t.saved = append(t.saved, metric)
	return t.innerSerializer.Serialize(metric)
}

func (t *testSerializer) SerializeBatch(metrics []telegraf.Metric) ([]byte, error) {
	t.saved = append(t.saved, metrics...)
	return t.innerSerializer.SerializeBatch(metrics)
}
