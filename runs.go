package rxengine

import (
	"fmt"
	"math/rand"
	"time"
)

//ExecuteRuns - runs the machine the specified number of times
func ExecuteRuns(numRuns float64, transitionMap *map[string][]Transition ) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if numRuns == 1 {
		timeTaken := singleRun(r, true, transitionMap)
		fmt.Printf("\n\nTime taken to be shipped = %v\n\n", timeTaken)
		return
	}

	totalDelay := 0.
	for index := 0.; index < numRuns; index++ {
		totalDelay += singleRun(r, false, transitionMap)
	}
	avgDelay := totalDelay / numRuns
	fmt.Printf("\n\nAvg time taken to be shipped for %v runs: %v \n\n", numRuns, avgDelay)
}

func singleRun(r *rand.Rand, isDemo bool, transitionMap *map[string][]Transition) float64 {
	currentDay := 0.
	rxMachine := NewRxMachine(currentDay, &StateTransitionMap)
	for rxMachine.FSM.Current() != DeliveredState {
		prob := r.Float64()
		currState := rxMachine.FSM.Current()
		var nextEvent Transition
		runningProb := 0.
		for _, v := range StateTransitionMap[currState] {
			if prob <= (runningProb + v.Prob) {
				nextEvent = v
				break
			}
			runningProb = runningProb + v.Prob
		}
		rxMachine.FSM.Event(nextEvent.EventType, nextEvent.Delay, isDemo)
	}

	return rxMachine.CurrentDay
}
