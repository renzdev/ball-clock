package ballclock

import "testing"

func TestBallClock(t *testing.T) {
	testBalls := 20
	result := RunSim(20, 0)
	if result {
		t.Errorf("Simulation failed for %d balls with %t result vs expected %t.", testBalls, result, false)
	}
}
