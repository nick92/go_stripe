package stripe

import (
	"net/http"
	"os"
	"trucknav/models"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentmethod"
	"github.com/stripe/stripe-go/v72/product"
	"github.com/stripe/stripe-go/v72/sub"
)

func addCustomer(c *gin.Context) {
	var addCustomerRequest models.AddCustomerRequest

	if err := c.ShouldBindJSON(&addCustomerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stripe.Key = os.Getenv("STRIPE_KEY")

	params := &stripe.CustomerParams{
		Name:  &addCustomerRequest.Name,
		Email: &addCustomerRequest.Email,
	}

	cust, err := customer.New(params)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payparams := &stripe.PaymentMethodAttachParams{
		Customer: &cust.ID,
	}

	pm, err := paymentmethod.Attach(
		addCustomerRequest.PaymentMethodId,
		payparams,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invoiceParams := &stripe.CustomerInvoiceSettingsParams{
		DefaultPaymentMethod: &pm.ID,
	}

	custparams := &stripe.CustomerParams{
		InvoiceSettings: invoiceParams,
	}

	_, err = customer.Update(
		cust.ID,
		custparams,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	prodparams := &stripe.ProductParams{}

	price, err := product.Get(addCustomerRequest.ProdId, prodparams)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subparams := &stripe.SubscriptionParams{
		Customer: stripe.String(cust.ID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: &price.DefaultPrice.ID,
			},
		},
	}

	s, err := sub.New(subparams)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := &models.AddCustomerResponse{
		Complete:       true,
		SubscriptionId: s.ID,
		Error:          "",
		CustomerId:     cust.ID,
	}

	c.JSON(http.StatusOK, gin.H{"complete": resp})
}

func getCustomer(c *gin.Context) {
	var customerId = c.Query("customer_id")

	stripe.Key = os.Getenv("STRIPE_KEY")

	params := &stripe.CustomerParams{}

	cust, err := customer.Get(customerId, params)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var subs []*stripe.Subscription

	subparams := &stripe.SubscriptionListParams{
		Customer: cust.ID,
	}
	i := sub.List(subparams)
	for i.Next() {
		s := i.Subscription()
		subs = append(subs, s)
	}

	pm, _ := paymentmethod.Get(
		cust.InvoiceSettings.DefaultPaymentMethod.ID,
		nil,
	)

	cust.InvoiceSettings.DefaultPaymentMethod = pm

	resp := &models.CustomerDetailsResponse{
		Complete:     true,
		Error:        "",
		Customer:     cust,
		Subscription: subs[0],
	}

	c.JSON(http.StatusOK, gin.H{"complete": resp})
}

func deleteCustomer(c *gin.Context) {

}

func createPaymentMethod(c *gin.Context) {
	stripe.Key = os.Getenv("STRIPE_KEY")

	params := &stripe.PaymentMethodParams{
		Card: &stripe.PaymentMethodCardParams{
			Number:   stripe.String("4242424242424242"),
			ExpMonth: stripe.String("8"),
			ExpYear:  stripe.String("2026"),
			CVC:      stripe.String("314"),
		},
		Type: stripe.String("card"),
	}
	pm, _ := paymentmethod.New(params)

	c.JSON(http.StatusOK, gin.H{"complete": pm.ID})
}
