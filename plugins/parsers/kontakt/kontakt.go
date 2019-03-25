package kontakt

import (
	"encoding/json"
	"fmt"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"time"
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
		return parseV3(metrics, json)
	case 4:
		return parseV4(metrics, json)
	default:
		return metrics, nil
	}
}

func parseV4(metrics []telegraf.Metric, json map[string]interface{}) ([]telegraf.Metric, error) {
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
				"sourceId": sourceId,
				"deviceAddress": address,
			},
			map[string]interface{}{
				"rssi": evt["rssi"],
				"data": bleEvt["data"],
				"srData": bleEvt["srData"],
			},
			time.Now(),
		)
		metrics = append(metrics, m)
	}


	return metrics, nil
}

func parseV3(metrics []telegraf.Metric, json map[string]interface{}) ([]telegraf.Metric, error) {
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
				"rssi": evt["rssi"],
				"data": evt["data"],
				"srData": evt["srData"],
				"sourceId": sourceId,
			},
			time.Now(),
		)
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