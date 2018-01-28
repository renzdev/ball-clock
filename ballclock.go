package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"
)

func main() {
	ballPtr := flag.Int("balls", 0, "ball count for simulation")
	limitPtr := flag.Int("limit", 0, "(optional) time limit for simulation specified in minutes")
	flag.Parse()

	// no real use for function return values here since RunSim outputs the results to console already
	// probably an indication of some output redundancy that could be cleaned up
	RunSim(*ballPtr, *limitPtr)
}

const ballMin = 27
const ballMax = 127
const minuteSize = 4
const fMinuteSize = 11
const hourSize = 11

// RunSim runs a ballClock simulation according to completion requirements specified
// by input parameters and returns a success/failure bool and result string to
// facilitate (very) basic testing and benchmarking
func RunSim(ballCount, timeLimit int) (bool, string) {
	if ballCount < ballMin || ballCount > ballMax {
		// we would handle errors more carefully in production code
		result := "Error - invalid ballCount specified for simulation"
		fmt.Println(result)
		return false, result
	}

	// this case wasn't explicitly stated in specification but seems safe to avoid it
	if timeLimit < 0 {
		result := "Error - invalid timeLimit specified for simulation"
		fmt.Println(result)
		return false, result
	}

	// init ballClock fields
	// considered going with more OO design to enapsulate init logic as
	// a ballClock method but it didn't really seem like golang style?
	bc := new(ballClock)
	bc.Main = make([]int, ballCount)
	for i := 0; i < len(bc.Main); i++ {
		bc.Main[i] = i + 1
	}
	bc.Min = make([]int, 0, minuteSize)
	bc.FiveMin = make([]int, 0, fMinuteSize)
	bc.Hour = make([]int, 0, hourSize)

	// configure simulation according to input parameters specifying mode
	// using anon functions for completion condition / results reporting rather than
	// checking conditionals for "mode" repeatedly in function body. this might not be
	// considered good/safe golang style either, will need to study with more time
	var isComplete func() bool
	var reportResults func() string
	var minutesElapsed int
	fmt.Print("BallClock simulation configured")
	if timeLimit > 0 {
		fmt.Print(" for Mode 2 (Clock State)\n")
		isComplete = func() bool {
			return minutesElapsed == timeLimit
		}
		reportResults = func() string {
			return bc.String()
		}
	} else {
		fmt.Print(" for Mode 1 (Cycle Days)\n")
		isComplete = func() bool {
			// Mode 1 is checking for all balls returned to initial ordering
			// This requires all balls to be returned to Main slice which only
			// occurs on the hour. This means we only need to perform the
			// expensive comparison loop once every 60 ticks. In my testing this
			// is by far the worst bottleneck in the problem and there are
			// undoubtly more clever tricks that could make it even faster.

			// Since we know that the initial balls were generated with values
			// in ascending order, test comparison loop for Main[i] < Main[i+1]
			// rather than taking up extra space for storage of original data

			// Note: 	Testing shows significant performance drops for
			//			95, 113, 119, 123, 126 with this algorithm
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
		reportResults = func() string {
			return fmt.Sprintf("%d balls cycle after %d days.", ballCount, minutesElapsed/60/24)
		}
	}

	// run simulation and report/return results
	startTime := time.Now()
	for !isComplete() {
		bc.Tick()
		minutesElapsed++
	}
	simDuration := time.Since(startTime).Seconds()

	result := reportResults()
	fmt.Println(result)
	fmt.Printf("Completed in %d milliseconds (%f.3 seconds)\n", int(simDuration*1000), simDuration)
	return true, result
}

type ballClock struct {
	Min     []int
	FiveMin []int
	Hour    []int
	Main    []int
}

func (bc *ballClock) Tick() {
	// pop front ball from the slice "queue"
	var newBall int
	newBall, bc.Main = bc.Main[0], bc.Main[1:]

	// minute track
	if len(bc.Min) < minuteSize {
		bc.Min = append(bc.Min, newBall)
		return
	}

	for i := minuteSize - 1; i >= 0; i-- {
		bc.Main = append(bc.Main, bc.Min[i])
	}
	bc.Min = make([]int, 0, minuteSize)

	// five minute track
	if len(bc.FiveMin) < fMinuteSize {
		bc.FiveMin = append(bc.FiveMin, newBall)
		return
	}

	for i := fMinuteSize - 1; i >= 0; i-- {
		bc.Main = append(bc.Main, bc.FiveMin[i])
	}
	bc.FiveMin = make([]int, 0, fMinuteSize)

	// hour track
	if len(bc.Hour) < hourSize {
		bc.Hour = append(bc.Hour, newBall)
		return
	}

	for i := hourSize - 1; i >= 0; i-- {
		bc.Main = append(bc.Main, bc.Hour[i])
	}
	bc.Hour = make([]int, 0, hourSize)

	// return home little explorer ball
	bc.Main = append(bc.Main, newBall)
}

func (bc *ballClock) String() string {
	bcState, _ := json.Marshal(bc)
	return string(bcState)
}
