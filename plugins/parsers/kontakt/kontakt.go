package kontakt

import (
	"encoding/json"
	"fmt"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"math"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	millisInSecond = int64(time.Second / time.Millisecond)
)

var (
	eventsParsed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "telegraf_parsers_kontakt_parsed_events",
		Help: "Number of events in parsed requests",
	})
)

type KontaktEventParser struct {
	DefaultTags map[string]string
}

func (p *KontaktEventParser) parseObject(metrics []telegraf.Metric, json map[string]interface{}) ([]telegraf.Metric, error) {
	version, ok := json["version"].(float64)
	if !ok {
		return metrics, nil
	}
	switch version {
	case 3:
		return p.parseV3(metrics, json)
	case 4:
		return p.parseV4(metrics, json)
	default:
		return metrics, nil
	}
}

func (p *KontaktEventParser) parseV4(metrics []telegraf.Metric, json map[string]interface{}) ([]telegraf.Metric, error) {
	events, ok := json["events"].([]interface{})
	if !ok {
		return metrics, nil
	}
	sourceId, ok := json["sourceId"].(string)
	if !ok {
		return metrics, nil
	}

	for _, event := range events {
		evt, ok := event.(map[string]interface{})
		if !ok {
			continue
		}
		bleEvt, ok := evt["ble"].(map[string]interface{})
		if !ok {
			continue
		}
		address, ok := bleEvt["deviceAddress"].(string)
		if !ok {
			continue
		}

		m, _ := metric.New(
			"telemetry",
			map[string]string{
				"sourceId":      sourceId,
				"deviceAddress": address,
			},
			map[string]interface{}{
				"rssi":   evt["rssi"],
				"data":   bleEvt["data"],
				"srData": bleEvt["srData"],
			},
			time.Now(),
		)

		timestamp, ok := evt["timestamp"].(float64)
		if ok {
			timestampInt := int64(timestamp)
			m.AddField("gatewayTimestamp", p.normalizeTimestamp(timestampInt))
		}

		metrics = append(metrics, m)
	}

	return metrics, nil
}

func (p *KontaktEventParser) parseV3(metrics []telegraf.Metric, json map[string]interface{}) ([]telegraf.Metric, error) {
	events, ok := json["events"].([]interface{})
	if !ok {
		return metrics, nil
	}
	sourceId, ok := json["sourceId"].(string)
	if !ok {
		return metrics, nil
	}

	for _, event := range events {
		evt, ok := event.(map[string]interface{})
		if !ok {
			continue
		}
		address, ok := evt["deviceAddress"].(string)
		if !ok {
			continue
		}

		m, _ := metric.New(
			"telemetry",
			map[string]string{
				"deviceAddress": address,
			},
			map[string]interface{}{
				"rssi":     evt["rssi"],
				"data":     evt["data"],
				"srData":   evt["srData"],
				"sourceId": sourceId,
			},
			time.Now(),
		)

		timestamp, ok := evt["timestamp"].(float64)
		if ok {
			timestampInt := int64(timestamp)
			m.AddField("gatewayTimestamp", p.normalizeTimestamp(timestampInt))
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
