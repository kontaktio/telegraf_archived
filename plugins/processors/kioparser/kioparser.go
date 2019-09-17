package kioparser

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/processors"
	"math"
	"strings"
)

const blockHeaderLength = 2
const serviceDataByte byte = 0x16

const manufacturerDataByte byte = 0xFF
const ibeaconPayloadLength = 21
const ibeaconIdentifierLength = 20

const embeddedKontaktUUIDPrefix = 0x6A
const quuppaBeaconUUIDPrefix = 0x1A

var ibeaconPreamble = []byte{0x4C, 0x00, 0x02, 0x15}[:]
var quuppaConstantBytes = []byte{0x67, 0xF7, 0xDB, 0x34, 0xC4, 0x03, 0x8E, 0x5C, 0x0B, 0xAA, 0x97, 0x30, 0x56}[:]

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
	16: "UNIVERSAL_TAG",
	17: "BRACELET_TAG_2",
	18: "MOBILE_DEVICE",
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
				fields, tags, success := parseData(b, rssiFloat)
				if success {
					result = append(result, metric)
					for k, v := range fields {
						metric.AddField(k, v)
					}
					for k, v := range tags {
						metric.AddTag(k, v)
					}
					metric.RemoveField("data")
				}
			}
		}
	}
	return result
}

func parseData(data []byte, rssi float64) (map[string]interface{}, map[string]string, bool) {
	fields := make(map[string]interface{})
	tags := make(map[string]string)
	success := false
	buffer := bytes.NewBuffer(data)
	for buffer.Len() >= blockHeaderLength {
		blockLen, _ := buffer.ReadByte()
		if int(blockLen) > buffer.Len() {
			return map[string]interface{}{}, map[string]string{}, success
		}
		block := bytes.NewBuffer(buffer.Next(int(blockLen)))
		blockType, _ := block.ReadByte()
		switch blockType {
		case manufacturerDataByte:
			if bytes.Compare(block.Next(len(ibeaconPreamble)), ibeaconPreamble) != 0 {
				continue
			}
			if nextByte, e := block.ReadByte(); e == nil {
				if err := block.UnreadByte(); err != nil {
					continue
				}
				switch nextByte {
				case quuppaBeaconUUIDPrefix:
					ibeaconIdentifierBlock := bytes.NewBuffer(block.Next(ibeaconPayloadLength))
					success = convertQuuppa(ibeaconIdentifierBlock, fields, tags)
				case embeddedKontaktUUIDPrefix:
					kioIdentifierBlock := bytes.NewBuffer(block.Next(ibeaconIdentifierLength))
					success = convertKontaktFrame(kioIdentifierBlock, fields, rssi)
					if uniqueId, exists := fields["uniqueId"]; success && exists {
						strUniqueId := uniqueId.(string)
						/*
							Same as for Quuppa - this frame will only be advertised by mobile phones, where mac address
							may not be constant, that's why we need more "stable" trackingId - uniqueId
						*/
						tags["trackingId"] = strUniqueId
					}
				}
			}
		case serviceDataByte:
			success = convertKontaktFrame(block, fields, rssi)
		default:
			continue
		}
	}

	return fields, tags, success
}

func convertKontaktFrame(block *bytes.Buffer, fields map[string]interface{}, rssi float64) bool {
	if bytes.Compare(block.Next(len(kioUuid)), kioUuid) != 0 {
		return false
	}
	blockIdentifier, _ := block.ReadByte()
	fields["packetType"] = int64(blockIdentifier)
	switch blockIdentifier {
	case plainIdentifier:
		convertPlain(block, fields, rssi)
		return true
	case telemetryIdentifier:
		convertTelemetry(block, fields)
		return true
	case locationIdentifier:
		convertLocation(block, fields, rssi)
		return true
	default:
		return false
	}
}

func convertQuuppa(buffer *bytes.Buffer, fields map[string]interface{}, tags map[string]string) bool {
	if buffer.Len() != ibeaconPayloadLength {
		return false
	}
	if _, err := buffer.ReadByte(); err != nil { // Quuppa prefix - 0x1A
		return false
	}
	quuppaTagId := hex.EncodeToString(buffer.Next(6))
	if _, err := buffer.ReadByte(); err != nil { // Quuppa checksum - now ignored
		return false
	}
	if bytes.Compare(buffer.Next(len(quuppaConstantBytes)), quuppaConstantBytes) != 0 {
		return false
	}
	fields["frameType"] = int64(manufacturerDataByte)
	fields["model"] = models[18]
	fields["uniqueId"] = quuppaTagId
	/*
		Parsing of Quuppa iBeacon packet is done for mobile applications advertising some identifier.
		Because we can be sure, that advertising mobile mac address will be constant, we need to use
		something else as an identifier (trackingId), that's why it is overridden
	*/
	tags["trackingId"] = quuppaTagId
	return true
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
			result["batteryLevel"] = asFloatUnsigned(buffer.ReadByte())
		case 0x02:
			result["sensitivity"] = asFloatUnsigned(buffer.ReadByte())
			result["x"] = asFloatSigned(buffer.ReadByte())
			result["y"] = asFloatSigned(buffer.ReadByte())
			result["z"] = asFloatSigned(buffer.ReadByte())
			result["secondsSinceDoubleTap"] = float64(binary.LittleEndian.Uint16(buffer.Next(2)))
			result["secondsSinceThreshold"] = float64(binary.LittleEndian.Uint16(buffer.Next(2)))
		case 0x05:
			result["lightLevel"] = asFloatUnsigned(buffer.ReadByte())
			result["temperature"] = asFloatSigned(buffer.ReadByte())
		case 0x06:
			result["sensitivity"] = asFloatUnsigned(buffer.ReadByte())
			result["x"] = asFloatSigned(buffer.ReadByte())
			result["y"] = asFloatSigned(buffer.ReadByte())
			result["z"] = asFloatSigned(buffer.ReadByte())
		case 0x07:
			result["secondsSinceThreshold"] = float64(binary.LittleEndian.Uint16(buffer.Next(2)))
		case 0x08:
			result["secondsSinceDoubleTap"] = float64(binary.LittleEndian.Uint16(buffer.Next(2)))
		case 0x09:
			result["secondsSinceTap"] = float64(binary.LittleEndian.Uint16(buffer.Next(2)))
		case 0x0A:
			result["lightLevel"] = asFloatUnsigned(buffer.ReadByte())
		case 0x0B:
			result["temperature"] = asFloatSigned(buffer.ReadByte())
		case 0x0C:
			result["batteryLevel"] = asFloatUnsigned(buffer.ReadByte())
		case 0x0D:
			result["secondsSinceClick"] = float64(binary.LittleEndian.Uint16(buffer.Next(2)))
		case 0x0E:
			result["secondsDoubleClick"] = float64(binary.LittleEndian.Uint16(buffer.Next(2)))
		case 0x0F:
			result["utcTimestamp"] = float64(binary.LittleEndian.Uint32(buffer.Next(4)))
		case 0x11:
			result["clickId"] = asFloatUnsigned(buffer.ReadByte())
			result["secondsSinceClick"] = float64(binary.LittleEndian.Uint16(buffer.Next(2)))
		case 0x12:
			result["humidity"] = asFloatUnsigned(buffer.ReadByte())
		case 0x13:
			result["temperature"] = asFixedPoint88Signed(buffer.Next(2))
		case 0x16:
			result["movementId"] = asFloatUnsigned(buffer.ReadByte())
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
	result["uniqueId"] = strings.Trim(buffer.String(), "\x00")
	result["frameType"] = int64(plainIdentifier)

	if distance, err := calculateDistance(rssi, txPowerFloat); err == nil {
		result["distance"] = distance
	}
}

func asFloatUnsigned(b byte, _ error) float64 {
	return float64(b)
}

func asFloatSigned(b byte, _ error) float64 {
	result := float64(b)
	if b >= 128 {
		return -256 + result
	} else {
		return result
	}
}

func asFixedPoint88Signed(b []byte) float64 {
	result := float64(binary.LittleEndian.Uint16(b)) / 256.0
	if result >= 128.0 {
		return result - 256.0
	} else {
		return result
	}
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
