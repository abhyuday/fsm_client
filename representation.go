package rxengine

//StateTransitionMap - maps the states and corresponding events
var StateTransitionMap = map[string][]Transition{
	InitialRxState: []Transition{
		{EventType: CheckRefillsEvent, Prob: 1, Delay: 0.1, TargetState: CheckRefillsState},
	},
	HasNoRefillsState: []Transition{
		{EventType: ContactMDEvent, Prob: 1, Delay: 1, TargetState: ContactedMDState},
	},
	ContactedMDState: []Transition{
		{EventType: MDApprovedRefillsEvent, Prob: 1, Delay: 2, TargetState: HasRefillsState},
	},
	HasRefillsState: []Transition{
		{EventType: StartInsuranceVerificationEvent, Prob: 1, Delay: 0.5, TargetState: InsuranceToBeVerifiedState},
	},
	RefillNotDueState: []Transition{
		{EventType: VerifyInsuranceSuccessEvent, Prob: 1, Delay: 0.3, TargetState: InsuranceVerificationSuccessState},
	},
	AuthorizePaymentSuccessState: []Transition{
		{EventType: StartStockCheckEvent, Prob: 1, Delay: 0.4, TargetState: StockToBeVerfiedState},
	},
	InsuranceVerificationSuccessState: []Transition{
		{EventType: SendCopayApprovalEvent, Prob: 1, Delay: 0.4, TargetState: CopayToBeApprovedState},
	},
	CheckRefillsState: []Transition{
		{EventType: HasRefillsEvent, Prob: 0.5, Delay: 0.1, TargetState: HasRefillsState},
		{EventType: HasNoRefillsEvent, Prob: 0.5, Delay: 1, TargetState: HasNoRefillsState},
	},
	InsuranceToBeVerifiedState: []Transition{
		{EventType: VerifyInsuranceSuccessEvent, Prob: 0.5, Delay: 0.4, TargetState: InsuranceVerificationSuccessState},
		{EventType: VerifyInsuranceFailEvent, Prob: 0.5, Delay: 2, TargetState: RefillNotDueState},
	},
	StockToBeVerfiedState: []Transition{
		{EventType: StockCheckFailEvent, Prob: 0.5, Delay: 2, TargetState: StockToBeVerfiedState},
		{EventType: StockCheckSuccessEvent, Prob: 0.5, Delay: 0.3, TargetState: PaymentToBeChargedState},
	},
	CopayToBeApprovedState: []Transition{
		{EventType: CopayApprovedEvent, Prob: 0.5, Delay: 0.3, TargetState: AuthorizePaymentSuccessState},
		{EventType: CopayNotApprovedEvent, Prob: 0.5, Delay: 1, TargetState: CopayToBeApprovedState},
	},
	PaymentToBeChargedState: []Transition{
		{EventType: PaymentExceptionEvent, Prob: 0.5, Delay: 2, TargetState: PaymentToBeChargedState},
		{EventType: PaymentSuccessEvent, Prob: 0.5, Delay: 0.1, TargetState: ToBeShippedState},
	},
	ToBeShippedState: []Transition{
		{EventType: ShipExceptionEvent, Prob: 0.5, Delay: 2, TargetState: ToBeShippedState},
		{EventType: ShipSuccessEvent, Prob: 0.5, Delay: 0.1, TargetState: ShippedState},
	},
	ShippedState: []Transition{
		{EventType: DeliverySuccessEvent, Prob: 1, Delay: 2, TargetState: DeliveredState},
	},
}
