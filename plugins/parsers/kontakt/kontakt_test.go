package kontakt

import (
	"fmt"
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
         "timestamp":1518421585,
         "deviceAddress":"f9:10:e1:40:d4:f1"
      }
   ]
}`)
	parser := KontaktEventParser{}
	metrics, err := parser.Parse(packetJson)
	if err != nil {
		t.Fatalf("Unable to parse json %v", err)
	}
	fmt.Printf("%v", metrics)

	if len(metrics) != 1 {
		t.Fatalf("Unexpected metrics length: %v", len(metrics))
	}
}
