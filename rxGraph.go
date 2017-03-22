package rxengine

import (
	"fmt"
	"math/rand"
	"time"
)

//RxGraph representation of the state machine
type RxGraph struct {
	InitialState       string                  `json:"initialState" bson:"initialState"`
	FinalState         string                  `json:"finalState" bson:"finalState"`
	StateTransitionMap map[string][]Transition `json:"stateTransitionMap" bson:"stateTransitionMap"`
	runningDelay       float64
	r                  *rand.Rand
	isDemo             bool
}

//Init - initializes the variables
func (rxGraph *RxGraph) Init(isDemo bool) {
	rxGraph.InitialState = InitialRxState
	rxGraph.FinalState = ShippedState
	rxGraph.StateTransitionMap = StateTransitionMap
	rxGraph.runningDelay = 0.
	rxGraph.isDemo = isDemo
}

//Execute - executes the engine specified number of times
func (rxGraph *RxGraph) Execute(numRuns float64) float64 {
	i := 1.
	totalDelay := 0.
	for i <= numRuns {
		rxGraph.SingleRun(i)
		totalDelay += rxGraph.runningDelay
		i++
	}
	delayToTarget := totalDelay / numRuns
	rxGraph.updateMovingDelayToTargetForAllStates(delayToTarget)
	return delayToTarget
}

//SingleRun - executes the engine one time
func (rxGraph *RxGraph) SingleRun(currentRun float64) {
	rxGraph.r = rand.New(rand.NewSource(time.Now().UnixNano()))
	rxGraph.runningDelay = 0
	currentState := rxGraph.InitialState
	finalState := rxGraph.FinalState
	for currentState != finalState {
		prob := rxGraph.r.Float64()
		var transition *Transition
		runningProb := 0.
		for k, t := range rxGraph.StateTransitionMap[currentState] {
			if prob <= (runningProb + t.Prob) {
				transition = &rxGraph.StateTransitionMap[currentState][k]
				break
			}
			runningProb = runningProb + t.Prob
		}

		if rxGraph.isDemo {
			fmt.Printf("\nTransition: %s -> %s\n", currentState, transition.TargetState)
			fmt.Println("Cuurent Delay: ", rxGraph.runningDelay, " Delay added: ", transition.Delay)
		}

		rxGraph.runningDelay = rxGraph.runningDelay + transition.Delay
		currentState = transition.TargetState
		transition.DelayFromOrgin = rxGraph.runningDelay
		transition.MovingDelayFromOrigin = movingAvg(transition.MovingDelayFromOrigin, rxGraph.runningDelay, currentRun)
	}
}

func (rxGraph *RxGraph) updateMovingDelayToTargetForAllStates(movingDelayFromTarget float64) {
	for k := range rxGraph.StateTransitionMap {
		for tIndex := range rxGraph.StateTransitionMap[k] {
			rxGraph.StateTransitionMap[k][tIndex].MovingDelayToTarget = movingDelayFromTarget - rxGraph.StateTransitionMap[k][tIndex].MovingDelayFromOrigin
		}
	}
}

func movingAvg(avg float64, newVal float64, total float64) float64 {
	return avg + (newVal-avg)/total
}
