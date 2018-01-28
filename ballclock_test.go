package ballclock

import (
	"testing"
)

var ballClockTests = []struct {
	ballCount    int
	timeLimit    int
	result       bool
	resultString string
}{
	{-1, 0, false, "Error - invalid ballCount specified for simulation"},
	{0, 0, false, "Error - invalid ballCount specified for simulation"},
	{26, 0, false, "Error - invalid ballCount specified for simulation"},
	{128, 0, false, "Error - invalid ballCount specified for simulation"},
	{30, 0, true, "30 balls cycle after 15 days."},
	{45, 0, true, "45 balls cycle after 378 days."},
	{30, 325, true, "{\"Min\":[],\"FiveMin\":[22,13,25,3,7],\"Hour\":[6,12,17,4,15],\"Main\":[11,5,26,18,2,30,19,8,24,10,29,20,16,21,28,1,23,14,27,9]}"},
}

func TestBallClockTable(t *testing.T) {
	for _, tt := range ballClockTests {
		result, resultString := RunSim(tt.ballCount, tt.timeLimit)
		if result != tt.result || resultString != tt.resultString {
			t.Errorf("Simulation test failed for ballCount %d, timeLimit %d.\n", tt.ballCount, tt.timeLimit)
			t.Errorf("Expected %t %s\n", tt.result, tt.resultString)
			t.Errorf("Got %t, %s\n", result, resultString)
		} else {
			t.Log(resultString)
		}

	}
}

func TestBallClockPerf(t *testing.T) {
	for i := 27; i <= 127; i++ {
		_, resultString := RunSim(i, 0)
		t.Log(resultString)
	}
}

func TestBallClockManual(t *testing.T) {
	_, resultString := RunSim(30, 325)
	t.Log(resultString)
}
