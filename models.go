package rxengine

import (
	"github.com/abhyuday/go-fsm"
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

//Transition - determines the probability that a particular event happens
type Transition struct {
	EventType string
	Prob      float64
	Delay     float64
    TargetState string
}

//RxMachine - determines the current state of the Rx as it flows throw different states
type RxMachine struct {
	CurrentDay float64
	FSM        *fsm.FSM
}
