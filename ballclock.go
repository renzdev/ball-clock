// update loop where each iteration represents one ball / minute
// loop operations:
//
// bottlenecks:	-checking for initial ball ordering on each loop
//				but only needs to occur once every 60 minutes?
//				give balls int val and check is only > comparison
//				-moving balls between 'tracks'?

package ballclock

import (
	"encoding/json"
	"fmt"
	"time"
)

// using const here for convenience but for production code where we might want more
// flexibility with simulation parameters it might be better to encapsulate
// these vars in a config struct?
const ballMin = 27
const ballMax = 127
const minuteSize = 4
const fMinuteSize = 11
const hourSize = 11

func RunSim(ballCount, minLimit int) bool {
	if ballCount <= ballMin && ballCount >= ballMin {
		// we would handle errors more carefully in production code
		fmt.Println("Error - no ballCount specified for simulation")
		return false
	}

	// init clock
	// could go further with OO design / ctor / etc for clock struct?
	bc := new(ballClock)
	bc.Min = make([]int, minuteSize)
	bc.FiveMin = make([]int, fMinuteSize)
	bc.Hour = make([]int, hourSize)
	bc.Main = make([]int, ballCount)
	for i := 0; i < len(bc.Main); i++ {
		bc.Main[i] = i + 1
	}

	// configure simulation
	// trying out anon functions rather than checking sim conditions repeatedly elsewhere
	var isComplete func() bool
	var reportResults func()
	var minutesElapsed int
	fmt.Print("BallClock simulation configured")
	if minLimit > 0 {
		fmt.Print(" for Mode 2 (Clock State)\n")
		isComplete = func() bool {
			return minutesElapsed == minLimit
		}
		reportResults = func() {
			fmt.Println(bc)
		}
	} else {
		fmt.Print(" for Mode 1 (Cycle Days)\n")
		isComplete = func() bool {
			if minutesElapsed > 0 && minutesElapsed%60 == 0 {
				for i := 0; i < len(bc.Main)-1; i++ {
					if bc.Main[i] >= bc.Main[i+1] {
						return false
					}
				}

				return true
			}
			return false
		}
		reportResults = func() {
			fmt.Printf("%d balls cycle after %d days.\n", ballCount, minutesElapsed*60/24)
		}
	}

	startTime := time.Now()

	for !isComplete() {
		bc.Tick()
		minutesElapsed++
	}

	simDuration := time.Since(startTime).Seconds()
	reportResults()
	fmt.Printf("Completed in %d milliseconds (%f.3 seconds)", int(simDuration*1000), simDuration)
	return true
}

type ballClock struct {
	Main    []int
	Min     []int
	FiveMin []int
	Hour    []int
}

func (bc *ballClock) Tick() {

}

func (bc *ballClock) String() string {
	bcState, _ := json.Marshal(bc)
	return string(bcState)
}
