package kontakt

import (
	"encoding/json"
	"fmt"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"github.com/influxdata/telegraf/plugins/parsers"
	"github.com/prometheus/client_golang/prometheus"
	"math"
	"time"
)

const (
	millisInSecond = int64(time.Second / time.Millisecond)
)

var eventsParsed = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "telegraf_parsers_kontakt_parsed_events",
	Help: "Number of events in parsed requests",
})

var gatewayToServerLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name:    "telegraf_gateway_to_server_latency_seconds",
	Help:    "Latency between gateway timestamp and server timestamp, measured in seconds.",
	Buckets: prometheus.ExponentialBuckets(0.032, 1.3, 24), // 32ms to ~17s
})

type KontaktEventParser struct {
	DefaultTags map[string]string
}

func (p *KontaktEventParser) parseObject(metrics []telegraf.Metric, json map[string]interface{}) ([]telegraf.Metric, error) {
	events, ok := json["events"].([]interface{})
	if !ok {
		return metrics, nil
	}
	sourceId, ok := json["sourceId"].(string)
	if !ok {
		return metrics, nil
	}

	sourceType, hasSourceType := json["sourceType"].(string)

	for _, event := range events {
		evt, ok := event.(map[string]interface{})
		if !ok {
			continue
		}
		address, ok := evt["deviceAddress"].(string)
		if !ok {
			continue
		}

		m := metric.New(
			"telemetry",
			map[string]string{
				"deviceAddress": address,
			},
			map[string]interface{}{
				"rssi":       evt["rssi"],
				"data":       evt["data"],
				"srData":     evt["srData"],
				"sourceId":   sourceId,
				"dataSource": "Kio",
			},
			time.Now(),
		)

		if hasSourceType {
			m.AddField("sourceType", sourceType)
		}

		timestamp, ok := evt["timestamp"].(float64)
		if ok {
			timestampInt := int64(timestamp)
			normalizedTimestamp := p.normalizeTimestamp(timestampInt)
			m.AddField("gatewayTimestamp", normalizedTimestamp)
			recordLatency(time.UnixMilli(normalizedTimestamp), time.Now())
		}

		metrics = append(metrics, m)
	}

	return metrics, nil
}

func (p *KontaktEventParser) Parse(buf []byte) ([]telegraf.Metric, error) {
	result := make([]telegraf.Metric, 0)
	var jsonOut map[string]interface{}

	err := json.Unmarshal(buf, &jsonOut)
	if err != nil {
		err = fmt.Errorf("unable to parse Kontakt Event, %s", err)
		return nil, err
	}
	metrics, err := p.parseObject(result, jsonOut)
	if err != nil {
		err = fmt.Errorf("unable to parse Kontakt Event, %s", err)
		return nil, err
	}

	metricsCount := len(metrics)
	eventsParsed.Add(float64(metricsCount))
	return metrics, nil
}

func (p *KontaktEventParser) ParseLine(line string) (telegraf.Metric, error) {
	metrics, err := p.Parse([]byte(line + "\n"))

	if err != nil {
		return nil, err
	}

	if len(metrics) < 1 {
		return nil, fmt.Errorf("can not parse the line: %s, for data format: json ", line)
	}

	return metrics[0], nil
}

func (p *KontaktEventParser) SetDefaultTags(tags map[string]string) {
	p.DefaultTags = tags
}

func (p *KontaktEventParser) normalizeTimestamp(timestamp int64) int64 {
	// Because of this comparison, it won't work for ms timestamps before 1970-01-25T20:31:23Z
	if timestamp > math.MaxInt32 {
		return timestamp
	} else {
		return timestamp * millisInSecond
	}
}

func recordLatency(gatewayTimestamp, serverTimestamp time.Time) {
	latency := serverTimestamp.Sub(gatewayTimestamp).Seconds()
	gatewayToServerLatency.Observe(latency)
}

func init() {
	parsers.Add("kontakt",
		func(defaultMetricName string) telegraf.Parser {
			return &KontaktEventParser{}
		})
	prometheus.MustRegister(eventsParsed)
	prometheus.MustRegister(gatewayToServerLatency)

}
