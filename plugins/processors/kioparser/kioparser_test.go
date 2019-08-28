package kioparser

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/metric"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func prepareMetric(data string) telegraf.Metric {
	var result, _ = metric.New(
		"NotImportant",
		map[string]string{
			"deviceAddress": "01:02:03:04:05:06",
		},
		map[string]interface{}{
			"data": data,
			"rssi": float64(-77),
		},
		time.Now())
	return result
}

func assertField(t *testing.T, metric telegraf.Metric, field string, value interface{}) {
	val, exists := metric.GetField(field)
	require.True(t, exists)
	require.Equal(t, value, val)
}

func assertTag(t *testing.T, metric telegraf.Metric, tag string, value string) {
	val, exists := metric.GetTag(tag)
	require.True(t, exists)
	require.Equal(t, value, val)
}

func TestParseLocation(t *testing.T) {
	parser := *New()

	var metric = prepareMetric("AgEGDhZq/gX0JgsARXNPT0lM")
	result := parser.Apply(metric)
	require.Equal(t, 1, len(result))
	assertField(t, result[0], "uniqueId", "EsOOIL")
	assertField(t, result[0], "moving", false)
	assertField(t, result[0], "channel", float64(38))
	assertField(t, result[0], "packetType", int64(locationIdentifier))
	require.True(t, result[0].HasField("distance"))

}

func TestParseEddystoneEID(t *testing.T) {
	parser := *New()

	var metric = prepareMetric("AgEGAwOq/g0Wqv4wEAECAwQFBgcI")
	result := parser.Apply(metric)
	require.Equal(t, 0, len(result))
}

func TestParseEmbeddedKontaktIBeacon(t *testing.T) {
	parser := *New()
	/*
		Generated using OVS:
		RawEventDataCreator.buildKontaktEmbededIBeaconAdvertise(
			"4.1",
			100,
			false,
			4,
			uniqueId,
			AdvertisedDeviceModel.MOBILE_DEVICE);
	*/
	var metric = prepareMetric("AgEGGv9MAAIVav4CEgQBZARhYmNkZWZnaGlqawAE")
	result := parser.Apply(metric)
	assertTag(t, result[0], "trackingId", "abcdefghijk")
	assertField(t, result[0], "uniqueId", "abcdefghijk")
}

func TestParseQuuppaIBeacon(t *testing.T) {
	parser := *New()
	/*
		Generated using OVS:
		buildQuuppaIBeaconAdvertise("00badbeef000")
	*/
	var metric = prepareMetric("AgEGGv9MAAIVGgC62+7wAABn99s0xAOOXAuqlzBW")
	result := parser.Apply(metric)

	assertTag(t, result[0], "trackingId", "00badbeef000")
	assertField(t, result[0], "uniqueId", "00badbeef000")
}

func TestParseOtherIBeacon(t *testing.T) {
	parser := *New()

	//Some iBeacon, not matching Kio/Quuppa criteria
	var metric = prepareMetric("AgEGGv9MAAIVzAC62+7wAABn99s0xAOOXAuqlzBW")

	result := parser.Apply(metric)
	require.Len(t, result, 0)
}

func TestParsePlain(t *testing.T) {
	parser := *New()
	/*
		Generated using OVS:
		RawEventDataCreator.buildNrf52BeaconAdvertise(
			"1.10",
			(byte) 100,
			false,
			(byte) 0x04,
			"abcdef",
			AdvertisedDeviceModel.BEACON_PRO)
	*/
	var metric1 = prepareMetric("DxZq/gIGAQpkBGFiY2RlZg==")
	result := parser.Apply(metric1)
	require.Equal(t, 1, len(result))
	assertField(t, result[0], "uniqueId", "abcdef")
	assertField(t, result[0], "model", "BEACON_PRO")
	require.False(t, result[0].HasField("data"))
	assertField(t, result[0], "packetType", int64(plainIdentifier))
	require.True(t, result[0].HasField("distance"))
}

func TestParseTelemetryNegativeValues(t *testing.T) {
	parser := *New()
	/*
		    Generated using OVS:
			RawEventDataCreator.buildTelemetryPacket(List.of(
		                new TemperatureField(-60))
	*/
	var metric1 = prepareMetric("BxZq/gMCC8Q=")

	/*
				Generated using OVS:
				RawEventDataCreator.buildTelemetryPacket(Set.of(
		                new Temperature16BitField(-60.25),
		                new AccelerationField(16, new byte[] {-10, -20, -30}))))
	*/
	var metric2 = prepareMetric("DhZq/gMFBhD27OIDE8DD")

	result1 := parser.Apply(metric1)
	parsedMetric := result1[0]
	assertField(t, parsedMetric, "temperature", float64(-60))

	result2 := parser.Apply(metric2)
	parsedMetric = result2[0]
	assertField(t, parsedMetric, "temperature", float64(-60.25))
	assertField(t, parsedMetric, "x", float64(-10))
	assertField(t, parsedMetric, "y", float64(-20))
	assertField(t, parsedMetric, "z", float64(-30))

}

func TestParseTelemetry(t *testing.T) {
	parser := *New()
	/*
			Generated using OVS:
			RawEventDataCreator.buildTelemetryPacket(List.of(
		                new TemperatureField(32),
		                new Temperature16BitField(32.125d),
		                new AccelerationField(16, new byte[] {1, 2, 3}),
		                new UTCTimeField(1010020030),
		                new HumidityField(32),
		                new DoubleTapField(32000),
		                new TapField(16000),
		                new IdentifiedMovementField(16, 3200),
		                new IdentifiedButtonClickField(55, 48000)
		        ))
	*/
	var metric1 = prepareMetric("LBZq/gMCCyADEyAgBQYQAQIDBQ++rjM8AhIgAwgAfQMJgD4EFhCADAQRN4C7")
	result := parser.Apply(metric1)
	require.Equal(t, 1, len(result))
	parsedMetric := result[0]
	assertField(t, parsedMetric, "temperature", 32.125)
	assertField(t, parsedMetric, "x", float64(1))
	assertField(t, parsedMetric, "y", float64(2))
	assertField(t, parsedMetric, "z", float64(3))
	assertField(t, parsedMetric, "sensitivity", float64(16))
	assertField(t, parsedMetric, "utcTimestamp", float64(1010020030))
	assertField(t, parsedMetric, "humidity", float64(32))
	assertField(t, parsedMetric, "secondsSinceDoubleTap", float64(32000))
	assertField(t, parsedMetric, "secondsSinceTap", float64(16000))
	assertField(t, parsedMetric, "movementId", float64(16))
	assertField(t, parsedMetric, "secondsSinceThreshold", float64(3200))
	assertField(t, parsedMetric, "clickId", float64(55))
	assertField(t, parsedMetric, "secondsSinceClick", float64(48000))
	assertField(t, result[0], "packetType", int64(telemetryIdentifier))
	require.False(t, result[0].HasField("distance"))

	require.False(t, parsedMetric.HasField("data"))
}
