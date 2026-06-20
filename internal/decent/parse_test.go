package decent

import "testing"

func TestParseMachine(t *testing.T) {
	got, err := ParseMachine([]byte(`{
		"timestamp":"2026-06-19T00:06:22.522380",
		"state":{"state":"sleeping","substate":"idle"},
		"flow":0.0,
		"pressure":0.00244140625,
		"targetFlow":0.0625,
		"targetPressure":0.0,
		"mixTemperature":41.57421875,
		"groupTemperature":42.07501220703125,
		"targetMixTemperature":20.0,
		"targetGroupTemperature":90.0,
		"profileFrame":4,
		"steamTemperature":103
	}`))
	if err != nil {
		t.Fatal(err)
	}
	if got.State != "sleeping" || got.Substate != "idle" {
		t.Fatalf("unexpected state: %#v", got)
	}
	if got.GroupTemperature != 42.07501220703125 {
		t.Fatalf("unexpected group temp: %v", got.GroupTemperature)
	}
	if got.Timestamp.IsZero() {
		t.Fatal("timestamp was not parsed")
	}
}
