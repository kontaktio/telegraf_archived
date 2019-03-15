package kioparser

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/processors"
)

const blockHeaderLength = 2
const serviceDataByte byte = 0x16

const plainIdentifier = 0x02
const minimumPlainLength = 9
const locationIdentifier byte = 0x05
const minimumLocationLength = 10

const telemetryIdentifier byte = 0x03

const externalPower byte = 0xFF

const movingFlagIdx = 0

var kioUuid = []byte{0x6A, 0xFE}[:]
var models = map[byte]string{
	1:  "SMART_BEACON",
	3:  "USB_BEACON",
	5:  "GATEWAY",
	4:  "CARD_BEACON",
	6:  "BEACON_PRO",
	8:  "TAG_BEACON",
	9:  "SMART_BEACON_3",
	10: "HEAVY_DUTY_BEACON",
	11: "CARD_BEACON_2",
	14: "TOUGH_BEACON_2",
	15: "BRACELET_TAG",
}

type KioParser struct {
}

var SampleConfig = ``

func (p *KioParser) SampleConfig() string {
	return SampleConfig
}

func (p *KioParser) Description() string {
	return "Parse Kontakt.io advertisement/scan response frames"
}

func (p *KioParser) Apply(metrics ...telegraf.Metric) []telegraf.Metric {
	for _, metric := range metrics {
		p.append("data", metric)
		metric.RemoveField("data")
		p.append("srData", metric)
		metric.RemoveField("srData")
	}
	return metrics
}

func (p *KioParser) append(fieldName string, metric telegraf.Metric) {
	field, exists := metric.GetField(fieldName)
	if exists {
		bytes, convError := base64.StdEncoding.DecodeString(field.(string))
		if convError == nil {
			fields := parseData(bytes)
			for k, v := range fields {
				metric.AddField(k, v)
			}
		}
	}
}

func parseData(data []byte) map[string]interface{} {
	result := make(map[string]interface{})
	buffer := bytes.NewBuffer(data)
	for buffer.Len() >= blockHeaderLength {
		blockLen, _ := buffer.ReadByte()
		if int(blockLen) > buffer.Len() {
			return make(map[string]interface{})
		}
		block := bytes.NewBuffer(buffer.Next(int(blockLen)))
		blockType, _ := block.ReadByte()
		if blockType != serviceDataByte {
			continue
		}
		if bytes.Compare(block.Next(len(kioUuid)), kioUuid) != 0 {
			continue
		}
		blockIdentifier, _ := block.ReadByte()
		switch blockIdentifier {
		case plainIdentifier:
			convertPlain(block, result)
		case telemetryIdentifier:
			convertTelemetry(block, result)
		case locationIdentifier:
			convertLocation(block, result)
		default:
			continue
		}
	}
	return result
}

func convertLocation(buffer *bytes.Buffer, result map[string]interface{}) {
	if buffer.Len() < minimumLocationLength {
		return
	}

	txPower, _ := buffer.ReadByte()
	result["txPower"] = float64(txPower)
	bleChannel, _ := buffer.ReadByte()
	result["channel"] = float64(bleChannel)
	model, _ := buffer.ReadByte()
	result["model"] = models[model]
	flags, _ := buffer.ReadByte()
	result["moving"] = flags&(1<<movingFlagIdx) == 1
	result["uniqueId"] = buffer.String()
}

func convertTelemetry(buffer *bytes.Buffer, result map[string]interface{}) {
	for buffer.Len() > 0 {
		fieldLength, _ := buffer.ReadByte()
		if fieldLength < 0 || buffer.Len() < int(fieldLength) {
			return
		}

		identifier, _ := buffer.ReadByte()
		switch identifier {
		case 0x01:
			result["utcTimestamp"] = float64(binary.LittleEndian.Uint32(buffer.Next(4)))
			result["batteryLevel"] = asFloat(buffer.ReadByte())
		case 0x02:
			result["sensitivity"] = asFloat(buffer.ReadByte())
			result["x"] = asFloat(buffer.ReadByte())
			result["y"] = asFloat(buffer.ReadByte())
			result["z"] = asFloat(buffer.ReadByte())
			result["secondsSinceDoubleTap"] = float64(binary.LittleEndian.Uint16(buffer.Next(2)))
			result["secondsSinceThreshold"] = float64(binary.LittleEndian.Uint16(buffer.Next(2)))
		case 0x05:
			result["lightLevel"] = asFloat(buffer.ReadByte())
			result["temperature"] = asFloat(buffer.ReadByte())
		case 0x06:
			result["sensitivity"] = asFloat(buffer.ReadByte())
			result["x"] = asFloat(buffer.ReadByte())
			result["y"] = asFloat(buffer.ReadByte())
			result["z"] = asFloat(buffer.ReadByte())
		case 0x07:
			result["secondsSinceThreshold"] = float64(binary.LittleEndian.Uint16(buffer.Next(2)))
		case 0x08:
			result["secondsSinceDoubleTap"] = float64(binary.LittleEndian.Uint16(buffer.Next(2)))
		case 0x09:
			result["secondsSinceTap"] = float64(binary.LittleEndian.Uint16(buffer.Next(2)))
		case 0x0A:
			result["lightLevel"] = asFloat(buffer.ReadByte())
		case 0x0B:
			result["temperature"] = asFloat(buffer.ReadByte())
		case 0x0C:
			result["batteryLevel"] = asFloat(buffer.ReadByte())
		case 0x0D:
			result["secondsSinceClick"] = float64(binary.LittleEndian.Uint16(buffer.Next(2)))
		case 0x0E:
			result["secondsDoubleClick"] = float64(binary.LittleEndian.Uint16(buffer.Next(2)))
		case 0x0F:
			result["utcTimestamp"] = float64(binary.LittleEndian.Uint32(buffer.Next(4)))
		case 0x11:
			result["clickId"] = asFloat(buffer.ReadByte())
			result["secondsSinceClick"] = float64(binary.LittleEndian.Uint16(buffer.Next(2)))
		case 0x12:
			result["humidity"] = asFloat(buffer.ReadByte())
		case 0x13:
			temperature16b := buffer.Next(2)
			result["temperature"] = float64(temperature16b[1]) + float64(float64(temperature16b[0])/256.0)
		case 0x16:
			result["movementId"] = asFloat(buffer.ReadByte())
			result["secondsSinceThreshold"] = float64(binary.LittleEndian.Uint16(buffer.Next(2)))
		default:
			buffer.Next(int(fieldLength) - 1)
		}
	}
}

func asFloat(b byte, _ error) float64 {
	return float64(b)
}

func convertPlain(buffer *bytes.Buffer, result map[string]interface{}) {
	if buffer.Len() < minimumPlainLength {
		return
	}

	model, _ := buffer.ReadByte()
	result["model"] = models[model]
	buffer.Next(2) //Skip firmware
	batteryLevel, _ := buffer.ReadByte()
	if batteryLevel != externalPower {
		result["batteryLevel"] = float64(batteryLevel)
	}
	txPower, _ := buffer.ReadByte()
	result["txPower"] = float64(txPower)
	result["uniqueId"] = buffer.String()
}

func New() *KioParser {
	return &KioParser{}
}

func init() {
	processors.Add("kioparser", func() telegraf.Processor {
		return &KioParser{}
	})
}
