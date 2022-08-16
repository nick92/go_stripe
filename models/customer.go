package models

import (
	"github.com/stripe/stripe-go/v72"
)

type AddCustomerRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	AddressLine1    string `json:"line1"`
	AddressLine2    string `json:"line2"`
	City            string `json:"city"`
	PostCode        string `json:"postcode"`
	Country         string `json:"country"`
	PaymentMethodId string `json:"payment_method_id"`
	ProdId          string `json:"product_id"`
}

type AddCustomerResponse struct {
	SubscriptionId string `json:"subscription_id"`
	CustomerId     string `json:"customer_id"`
	Complete       bool   `json:"complete"`
	Error          string `json:"error"`
}

type GetCustomerRequest struct {
	CustomerId string `json:"customer_id"`
}

type CustomerDetailsResponse struct {
	Customer     *stripe.Customer     `json:"customer"`
	Subscription *stripe.Subscription `json:"subscription"`
	Complete     bool                 `json:"complete"`
	Error        string               `json:"error"`
}

type CustomerSubscription struct {
}
