package kioparser

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/processors"
	"math"
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

var unknownTxPowerError = errors.New("TxPower unknown")

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

var rssisAt1M = map[int]float64{
	4:   -59,
	0:   -65,
	-4:  -69,
	-8:  -72,
	-12: -77,
	-16: -81,
	-20: -84,
	-30: -115,
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
	result := make([]telegraf.Metric, 0)
	for _, metric := range metrics {
		rssi, exists := metric.GetField("rssi")
		if !exists {
			continue
		}
		rssiFloat, ok := rssi.(float64)
		if !ok {
			continue
		}
		field, exists := metric.GetField("data")
		if exists {
			b, convError := base64.StdEncoding.DecodeString(field.(string))
			if convError == nil {
				fields, success := parseData(b, rssiFloat)
				if success {
					result = append(result, metric)
					for k, v := range fields {
						metric.AddField(k, v)
					}
					metric.RemoveField("data")
				}
			}
		}
	}
	return result
}

func parseData(data []byte, rssi float64) (map[string]interface{}, bool) {
	result := make(map[string]interface{})
	success := false
	buffer := bytes.NewBuffer(data)
	for buffer.Len() >= blockHeaderLength {
		blockLen, _ := buffer.ReadByte()
		if int(blockLen) > buffer.Len() {
			return make(map[string]interface{}), success
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
		result["packetType"] = int64(blockIdentifier)
		switch blockIdentifier {
		case plainIdentifier:
			success = true
			convertPlain(block, result, rssi)
		case telemetryIdentifier:
			success = true
			convertTelemetry(block, result)
		case locationIdentifier:
			success = true
			convertLocation(block, result, rssi)
		default:
			continue
		}
	}

	return result, success
}

func convertLocation(buffer *bytes.Buffer, result map[string]interface{}, rssi float64) {
	if buffer.Len() < minimumLocationLength {
		return
	}
	result["frameType"] = int64(locationIdentifier)

	txPower, _ := buffer.ReadByte()
	txPowerFloat := float64(txPower)
	result["txPower"] = txPowerFloat
	bleChannel, _ := buffer.ReadByte()
	result["channel"] = float64(bleChannel)
	model, _ := buffer.ReadByte()
	result["model"] = models[model]
	flags, _ := buffer.ReadByte()
	result["moving"] = flags&(1<<movingFlagIdx) == 1
	result["uniqueId"] = buffer.String()

	if distance, err := calculateDistance(rssi, txPowerFloat); err == nil {
		result["distance"] = distance
	}

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
	result["frameType"] = int64(telemetryIdentifier)
}

func convertPlain(buffer *bytes.Buffer, result map[string]interface{}, rssi float64) {
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
	txPowerFloat := float64(txPower)
	result["txPower"] = txPowerFloat
	result["uniqueId"] = buffer.String()
	result["frameType"] = int64(plainIdentifier)

	if distance, err := calculateDistance(rssi, txPowerFloat); err == nil {
		result["distance"] = distance
	}
}

func asFloat(b byte, _ error) float64 {
	return float64(b)
}

func calculateDistance(rssi, txPower float64) (float64, error) {
	var txPowerNormalized int
	if txPower <= 4 { //Values 0 and 4
		txPowerNormalized = int(txPower)
	} else {
		txPowerNormalized = int(txPower) - 256
	}
	rssiAt1M, exists := rssisAt1M[txPowerNormalized]
	if !exists {
		return 0, unknownTxPowerError
	}

	ratio := 1.0 * rssi / rssiAt1M
	if ratio < 1.0 {
		return math.Pow(ratio, 10), nil
	} else {
		return 0.89976*math.Pow(ratio, 7.7095) + 0.111, nil
	}
}

func New() *KioParser {
	return &KioParser{}
}

func init() {
	processors.Add("kioparser", func() telegraf.Processor {
		return &KioParser{}
	})
}
