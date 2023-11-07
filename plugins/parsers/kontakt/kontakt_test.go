package kontakt

import (
	"github.com/influxdata/telegraf"
	"testing"
)

func TestParsePacket(t *testing.T) {
	packetJson := []byte(`
{
   "version":3,
   "timestamp":1518421586,
   "sourceId":"gatewayUniqueId",
   "sourceType":"GATEWAY",
   "events":[
      {
         "rssi":-51,
         "data":"DhZq/gIIBAEBBHlhdW0=",
         "srData":"someSrDataThatShouldBeBase64",
         "timestamp":1688335200,
         "deviceAddress":"f9:10:e1:40:d4:f1"
      }
   ]
}`)
	parser := KontaktEventParser{}
	metrics, err := parser.Parse(packetJson)
	if err != nil {
		t.Fatalf("Unable to parse json %v", err)
	}

	if len(metrics) != 1 {
		t.Fatalf("Unexpected metrics length: %v", len(metrics))
	}

	metric := metrics[0]
	deviceAddress, ok := metric.GetTag("deviceAddress")
	if !ok || deviceAddress != "f9:10:e1:40:d4:f1" {
		t.Fatalf("missing or incorrect deviceAddress: %v", deviceAddress)
	}

	AssertHasField(t, metric, "rssi", float64(-51))
	AssertHasField(t, metric, "data", "DhZq/gIIBAEBBHlhdW0=")
	AssertHasField(t, metric, "srData", "someSrDataThatShouldBeBase64")
	AssertHasField(t, metric, "sourceId", "gatewayUniqueId")
	AssertHasField(t, metric, "sourceType", "GATEWAY")
	AssertHasField(t, metric, "gatewayTimestamp", int64(1688335200000))
	AssertHasField(t, metric, "dataSource", "Kio")

}

func AssertHasField(t *testing.T, metric telegraf.Metric, fieldName string, expectedValue interface{}) {
	fieldValue, ok := metric.GetField(fieldName)
	if !ok || fieldValue != expectedValue {
		t.Fatalf("Missing or incorrect %v. Expected %v, got %v", fieldName, expectedValue, fieldValue)
	}
}
