package stripe

import (
	"net/http"
	"os"
	"trucknav/models"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentintent"
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
		resp := &models.CustomerDetailsResponse{
			Complete: false,
			Error:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, gin.H{"complete": resp})
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

	pm, err := paymentmethod.Get(
		cust.InvoiceSettings.DefaultPaymentMethod.ID,
		nil,
	)

	if err != nil {
		resp := &models.CustomerDetailsResponse{
			Complete: false,
			Error:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, gin.H{"complete": resp})
		return
	}

	cust.InvoiceSettings.DefaultPaymentMethod = pm

	var prod *stripe.Product
	prodparams := &stripe.ProductParams{}

	if len(subs) > 0 {
		prod, err = product.Get(
			subs[0].Plan.Product.ID,
			prodparams,
		)

		if err != nil {
			resp := &models.CustomerDetailsResponse{
				Complete: false,
				Error:    err.Error(),
			}
			c.JSON(http.StatusBadRequest, gin.H{"complete": resp})
			return
		}
	}

	cust.InvoiceSettings.DefaultPaymentMethod = pm

	var paymentIntents []*stripe.PaymentIntent
	piparams := &stripe.PaymentIntentListParams{
		Customer: &cust.ID,
	}
	piparams.Filters.AddFilter("limit", "", "5")
	pi := paymentintent.List(piparams)
	for pi.Next() {
		payment := pi.PaymentIntent()
		paymentIntents = append(paymentIntents, payment)
	}

	resp := &models.CustomerDetailsResponse{
		Complete:      true,
		Error:         "",
		Customer:      cust,
		PaymentMethod: pm,
		ActiveProduct: prod,
		Payments:      paymentIntents,
		Subscription:  subs,
	}

	c.JSON(http.StatusOK, gin.H{"complete": resp})
}

func deleteCustomer(c *gin.Context) {
	stripe.Key = os.Getenv("STRIPE_KEY")
	var customerId = c.Query("customer_id")

	_, err := customer.Del(customerId, nil)

	if err != nil {
		resp := &models.CustomerDetailsResponse{
			Complete: false,
			Error:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, gin.H{"complete": resp})
		return
	}

	c.JSON(http.StatusOK, gin.H{"complete": true})
}

func addCustomerSubscription(c *gin.Context) {
	stripe.Key = os.Getenv("STRIPE_KEY")
	var addCustomerSubRequest models.AddCustomerSubscriptionRequest

	if err := c.ShouldBindJSON(&addCustomerSubRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	custParams := &stripe.CustomerParams{}

	cust, err := customer.Get(addCustomerSubRequest.CustomerId, custParams)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if cust == nil {
		resp := &models.AddCustomerResponse{
			Complete:   false,
			Error:      "Missing or invalid customer acount",
			CustomerId: addCustomerSubRequest.CustomerId,
		}

		c.JSON(http.StatusBadRequest, gin.H{"complete": resp})
		return
	}

	prodparams := &stripe.ProductParams{}

	price, err := product.Get(addCustomerSubRequest.ProductId, prodparams)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subparams := &stripe.SubscriptionParams{
		Customer: stripe.String(addCustomerSubRequest.CustomerId),
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
		CustomerId:     addCustomerSubRequest.CustomerId,
	}

	c.JSON(http.StatusOK, gin.H{"complete": resp})
}

func cancelCustomerSubscription(c *gin.Context) {
	stripe.Key = os.Getenv("STRIPE_KEY")
	var subId = c.Query("subscription_id")

	_, err := sub.Cancel(
		subId,
		nil,
	)

	if err != nil {
		resp := &models.CustomerDetailsResponse{
			Complete: false,
			Error:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, gin.H{"complete": resp})
		return
	}

	c.JSON(http.StatusOK, gin.H{"complete": true})
}

func getCustomerPaymentMethods(c *gin.Context) {
	stripe.Key = os.Getenv("STRIPE_KEY")
	var customerId = c.Query("customer_id")
	var resp []*stripe.PaymentMethod

	custparams := &stripe.CustomerParams{}

	cust, err := customer.Get(customerId, custparams)

	if err != nil {
		resp := &models.CustomerDetailsResponse{
			Complete: false,
			Error:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, gin.H{"complete": resp})
		return
	}

	params := &stripe.PaymentMethodListParams{
		Customer: stripe.String(customerId),
		Type:     stripe.String("card"),
	}
	i := paymentmethod.List(params)

	for i.Next() {
		pm := i.PaymentMethod()
		if pm.ID == cust.InvoiceSettings.DefaultPaymentMethod.ID {
			pm.Card.Description = "default"
		}
		resp = append(resp, pm)
	}

	c.JSON(http.StatusOK, gin.H{"complete": resp})
}

func setCustomerDefaultPaymentMethods(c *gin.Context) {
	stripe.Key = os.Getenv("STRIPE_KEY")
	var customerId = c.Query("customer_id")
	var paymentId = c.Query("payment_id")

	custparams := &stripe.CustomerParams{}

	cust, err := customer.Get(customerId, custparams)

	if err != nil {
		resp := &models.CustomerDetailsResponse{
			Complete: false,
			Error:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, gin.H{"complete": resp})
		return
	}

	invoiceParams := &stripe.CustomerInvoiceSettingsParams{
		DefaultPaymentMethod: &paymentId,
	}

	ncustparams := &stripe.CustomerParams{
		InvoiceSettings: invoiceParams,
	}

	_, err = customer.Update(
		cust.ID,
		ncustparams,
	)

	if err != nil {
		resp := &models.CustomerDetailsResponse{
			Complete: false,
			Error:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, gin.H{"complete": resp})
		return
	}

	c.JSON(http.StatusOK, gin.H{"complete": true})
}

func createNewCustomerPaymentMethod(c *gin.Context) {
	stripe.Key = os.Getenv("STRIPE_KEY")

	var addCustomerPaymentRequest models.AddCustomerPaymentMethodRequest

	if err := c.ShouldBindJSON(&addCustomerPaymentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payparams := &stripe.PaymentMethodAttachParams{
		Customer: &addCustomerPaymentRequest.CustomerId,
	}

	pm, err := paymentmethod.Attach(
		addCustomerPaymentRequest.PaymentMethodId,
		payparams,
	)

	if err != nil {
		resp := &models.CustomerDetailsResponse{
			Complete: false,
			Error:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, gin.H{"complete": resp})
		return
	}

	c.JSON(http.StatusOK, gin.H{"complete": pm.ID})
}

func removeCustomerPaymentMethod(c *gin.Context) {
	stripe.Key = os.Getenv("STRIPE_KEY")
	var paymentId = c.Query("payment_id")

	payparams := &stripe.PaymentMethodDetachParams{}

	pm, err := paymentmethod.Detach(
		paymentId,
		payparams,
	)

	if err != nil {
		resp := &models.CustomerDetailsResponse{
			Complete: false,
			Error:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, gin.H{"complete": resp})
		return
	}

	c.JSON(http.StatusOK, gin.H{"complete": pm.ID})
}

func CreatePaymentMethod(key string) string {
	stripe.Key = key

	params := &stripe.PaymentMethodParams{
		Card: &stripe.PaymentMethodCardParams{
			Number:   stripe.String("4242424242424242"),
			ExpMonth: stripe.String("8"),
			ExpYear:  stripe.String("2030"),
			CVC:      stripe.String("314"),
		},
		Type: stripe.String("card"),
	}
	pm, _ := paymentmethod.New(params)

	return pm.ID
}
