package main

import (
	"fmt"
	"time"
	"github.com/abhyuday/go-fsm"
	"math/rand"
	"os"
	"strconv"
	"encoding/json"
)

//STATES
const (
	InitialRxState                    = "initail_rx_st"
	TransferCompletedState            = "tranfer_completed_st"
	CheckRefillsState                 = "check_refills_st"
	HasRefillsState                   = "has_refills_st"
	HasNoRefillsState                 = "has_no_refills_st"
	ContactedMDState                  = "contacted_md_st"
	InsuranceToBeVerifiedState        = "insurance_to_be_verified_st"
	InsuranceVerificationSuccessState = "insurance_verification_success_st"
	InsuranceVerificationFailState    = "insurance_verification_fail_st"
	RefillNotDueState                 = "refill_not_due_st"
	CopayToBeApprovedState            = "copay_approval_sent_st"
	AuthorizePaymentSuccessState      = "authorize_payment_success_st"
	StockToBeVerfiedState             = "stock_to_be_verified_st"
	PaymentToBeChargedState           = "payment_to_be_charged_st"
	ToBeShippedState                  = "to_be_shipped_st"
	ShippedState                      = "shipped_st"
	DeliveredState                    = "delivered_st"
	ArchivedState                     = "archived_st" //not relevant for refill authorization
)

//Events
const (
	CheckRefillsEvent               = "check_refills_evt"
	HasRefillsEvent                 = "has_refills_evt"
	HasNoRefillsEvent               = "has_no_refills_evt"
	ContactMDEvent                  = "contact_md_evt"
	MDApprovedRefillsEvent          = "md_approved_refills_evt"
	StartInsuranceVerificationEvent = "start_insurance_verification_evt"
	VerifyInsuranceSuccessEvent     = "verified_insurance_success_evt"
	VerifyInsuranceFailEvent        = "verified_insurance_fail_evt"
	SendCopayApprovalEvent          = "send_copay_approval_evt"
	CopayApprovedEvent              = "copay_approved_evt"
	CopayNotApprovedEvent           = "copay_not_approved_evt"
	AuthorizePaymentEvent           = "authorize_payment_evt"
	AuthorizePaymentSuccessEvent    = "authorize_payment_success_evt"
	StartStockCheckEvent            = "start_stock_check_evt"
	StockCheckSuccessEvent          = "stock_check_success_evt"
	StockCheckFailEvent             = "stock_check_fail_evt"
	StartPaymentChargeEvent         = "start_payment_charge_evt"
	PaymentExceptionEvent           = "payment_exception_evt"
	PaymentSuccessEvent             = "payment_success_evt"
	ShipExceptionEvent              = "ship_exception_evt"
	ShipSuccessEvent                = "ship_success_evt"
	DeliverySuccessEvent            = "delivery_sucess_evt"
)

//RxMachine - determines the current state of the Rx as it flows throw different states
type RxMachine struct {
	CurrentDay float64
	FSM        *fsm.FSM
}

//EventProbability - determines the probability that a particular event happens
type EventProbability struct {
	EventType string
	Prob      float64
	Delay     float64
}

var stateEventMap = map[string][]EventProbability{
	InitialRxState: []EventProbability{
		{EventType: CheckRefillsEvent, Prob: 1, Delay: 0.1},
	},
	HasNoRefillsState: []EventProbability{
		{EventType: ContactMDEvent, Prob: 1, Delay: 1},
	},
	ContactedMDState: []EventProbability{
		{EventType: MDApprovedRefillsEvent, Prob: 1, Delay: 2},
	},
	HasRefillsState: []EventProbability{
		{EventType: StartInsuranceVerificationEvent, Prob: 1, Delay: 0.5},
	},
	RefillNotDueState: []EventProbability{
		{EventType: VerifyInsuranceSuccessEvent, Prob: 1, Delay: 0.3},
	},
	AuthorizePaymentSuccessState: []EventProbability{
		{EventType: StartStockCheckEvent, Prob: 1, Delay: 0.4},
	},
	InsuranceVerificationSuccessState: []EventProbability{
		{EventType: SendCopayApprovalEvent, Prob: 1, Delay: 0.4},
	},
	CheckRefillsState: []EventProbability{
		{EventType: HasRefillsEvent, Prob: 0.5, Delay: 0.1},
		{EventType: HasNoRefillsEvent, Prob: 0.5, Delay: 1},
	},
	InsuranceToBeVerifiedState: []EventProbability{
		{EventType: VerifyInsuranceSuccessEvent, Prob: 0.5, Delay: 0.4},
		{EventType: VerifyInsuranceFailEvent, Prob: 0.5, Delay: 2},
	},
	StockToBeVerfiedState: []EventProbability{
		{EventType: StockCheckFailEvent, Prob: 0.5, Delay: 2},
		{EventType: StockCheckSuccessEvent, Prob: 0.5, Delay: 0.3},
	},
	CopayToBeApprovedState: []EventProbability{
		{EventType: CopayApprovedEvent, Prob: 0.5, Delay: 0.3},
		{EventType: CopayNotApprovedEvent, Prob: 0.5, Delay: 1},
	},
	PaymentToBeChargedState: []EventProbability{
		{EventType: PaymentExceptionEvent, Prob: 0.5, Delay: 2},
		{EventType: PaymentSuccessEvent, Prob: 0.5, Delay: 0.1},
	},
	ToBeShippedState: []EventProbability{
		{EventType: ShipExceptionEvent, Prob: 0.5, Delay: 2},
		{EventType: ShipSuccessEvent, Prob: 0.5, Delay: 0.1},
	},
	ShippedState: []EventProbability{
		{EventType: DeliverySuccessEvent, Prob: 1, Delay: 2},
	},
}

func testFunction()  {
	
}

var events = fsm.Events{

	//check refills flows
	{Name: CheckRefillsEvent, Src: []string{InitialRxState}, Dst: CheckRefillsState},
	{Name: HasRefillsEvent, Src: []string{CheckRefillsState}, Dst: HasRefillsState},
	{Name: HasNoRefillsEvent, Src: []string{CheckRefillsState}, Dst: HasNoRefillsState},
	{Name: ContactMDEvent, Src: []string{HasNoRefillsState}, Dst: ContactedMDState},
	{Name: MDApprovedRefillsEvent, Src: []string{ContactedMDState}, Dst: HasRefillsState},

	//verify insurance flow
	{Name: StartInsuranceVerificationEvent, Src: []string{HasRefillsState}, Dst: InsuranceToBeVerifiedState},
	{Name: VerifyInsuranceSuccessEvent, Src: []string{InsuranceToBeVerifiedState}, Dst: InsuranceVerificationSuccessState},
	{Name: VerifyInsuranceFailEvent, Src: []string{InsuranceToBeVerifiedState}, Dst: RefillNotDueState},
	{Name: VerifyInsuranceSuccessEvent, Src: []string{RefillNotDueState}, Dst: InsuranceVerificationSuccessState},

	//copay approval flow
	{Name: SendCopayApprovalEvent, Src: []string{InsuranceVerificationSuccessState}, Dst: CopayToBeApprovedState},
	{Name: CopayApprovedEvent, Src: []string{CopayToBeApprovedState}, Dst: AuthorizePaymentSuccessState},
	{Name: CopayNotApprovedEvent, Src: []string{CopayToBeApprovedState}, Dst: CopayToBeApprovedState},

	//stock check flow
	{Name: StartStockCheckEvent, Src: []string{AuthorizePaymentSuccessState}, Dst: StockToBeVerfiedState},
	{Name: StockCheckFailEvent, Src: []string{StockToBeVerfiedState}, Dst: StockToBeVerfiedState},
	{Name: StockCheckSuccessEvent, Src: []string{StockToBeVerfiedState}, Dst: PaymentToBeChargedState},

	// payment flow
	{Name: PaymentExceptionEvent, Src: []string{PaymentToBeChargedState}, Dst: PaymentToBeChargedState},
	{Name: PaymentSuccessEvent, Src: []string{PaymentToBeChargedState}, Dst: ToBeShippedState},

	//shipping flow
	{Name: ShipExceptionEvent, Src: []string{ToBeShippedState}, Dst: ToBeShippedState},
	{Name: ShipSuccessEvent, Src: []string{ToBeShippedState}, Dst: ShippedState},
	{Name: DeliverySuccessEvent, Src: []string{ShippedState}, Dst: DeliveredState},
}

//NewRxMachine - Creates the RxMachine instance
func NewRxMachine(currentDay float64) *RxMachine {
	d := &RxMachine{
		CurrentDay: currentDay,
	}

	d.FSM = fsm.NewFSM(
		InitialRxState,
		events,
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) { d.enterState(e) },
		},
	)
	return d
}

func (d *RxMachine) enterState(e *fsm.Event) {
	
	if e.Args == nil {
		return
	}

	isDemo := e.Args[1].(bool)
	if isDemo {
		fmt.Printf("\nTransition: %s -> %s\n", e.Src, e.Dst)
	}
	delay := e.Args[0].(float64)
	if isDemo {
 		fmt.Println("Day: ", d.CurrentDay, " Delay added: ", delay)
	}
	d.CurrentDay += delay
}

func singleRun(r *rand.Rand, isDemo bool) float64 {
	currentDay := 0.
	rxMachine := NewRxMachine(currentDay)
	for rxMachine.FSM.Current() != ShippedState {
		prob := r.Float64()
		currState := rxMachine.FSM.Current()
		var nextEvent EventProbability
		runningProb := 0.
		for _, v := range stateEventMap[currState] {
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

func executeRuns(numRuns float64)  {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if numRuns == 1 {
		timeTaken := singleRun(r, true)
		fmt.Printf("\n\nTime taken to be shipped = %v\n\n", timeTaken)
		return
	}

	totalDelay := 0.
	for index := 0.; index < numRuns; index++ {
		totalDelay += singleRun(r, false)
	}
	avgDelay := totalDelay/numRuns
	fmt.Printf("\n\nAvg time taken to be shipped for %v runs: %v \n\n", numRuns, avgDelay)
}
func main() {
	if len(os.Args) == 1 {
		fmt.Println("Please specify num of runs")
		return
	}
	numRuns, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Println("Please enter valid num of runs")
		return
	}
	executeRuns(numRuns)
	jsonString, err := json.MarshalIndent(stateEventMap, "", " ")
	f, err := os.Create("stateMap.txt")
	f.WriteString(string(jsonString))
}
