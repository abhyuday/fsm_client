package rxengine

import (
	"fmt"
	"github.com/abhyuday/go-fsm"
)

var rxMachineEvents = fsm.Events{

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
func NewRxMachine(currentDay float64, transitionMap *map[string][]Transition) *RxMachine {
	d := &RxMachine{
		CurrentDay: currentDay,
	}

	rxMachineEventList := rxMachineEventListFromTransitionMap(transitionMap)
	//fmt.Printf("%+v\n\n", rxMachineEventList)
	//fmt.Printf("%+v\n\n", rxMachineEvents)
	d.FSM = fsm.NewFSM(
		InitialRxState,
		rxMachineEventList,
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) { d.ProcessState(e) },
		},
	)
	return d
}

func rxMachineEventListFromTransitionMap(transitionMap *map[string][]Transition) fsm.Events {
	rxMachineEventList := fsm.Events{}
	for k, transitionArr := range *transitionMap {
		src := []string{k}
		for _, transition := range transitionArr {
			event := fsm.EventDesc{
				Name: transition.EventType,
				Dst: transition.TargetState,
				Src: src,
			}
			rxMachineEventList = append(rxMachineEventList, event)
		}
	}
	return rxMachineEventList
}
//ProcessState - does action as the state machine enters the state
func (d *RxMachine) ProcessState(e *fsm.Event) {
	if e.Args == nil {
		return
	}
	isDemo := e.Args[1].(bool)
	delay := e.Args[0].(float64)
	if isDemo {
		fmt.Printf("\nTransition: %s -> %s\n", e.Src, e.Dst)
		fmt.Println("Day: ", d.CurrentDay, " Delay added: ", delay)
	}
	d.CurrentDay += delay
}
