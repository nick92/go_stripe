package stripe

import (
	"net/http"
	"os"

	"trucknav/models"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentmethod"
	"github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/product"
	"github.com/stripe/stripe-go/v72/sub"
)

func InitCustomerRoutes(router *gin.RouterGroup) {
	router.GET("/get", getCustomer)
	router.POST("/add", addCustomer)
	router.DELETE("/delete", deleteCustomer)
}

func InitPriceRoutes(router *gin.RouterGroup) {
	router.GET("/get", getPrices)
}

func addCustomer(c *gin.Context) {
	var addCustomerRequest models.AddCustomerRequest

	if err := c.ShouldBindJSON(&addCustomerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stripe.Key = os.Getenv("STRIPE_KEY")

	params := &stripe.CustomerParams{
		Name: &addCustomerRequest.Name,
		// Address: &stripe.AddressParams{
		// 	Line1:      &addCustomerRequest.AddressLine1,
		// 	Line2:      &addCustomerRequest.AddressLine2,
		// 	PostalCode: &addCustomerRequest.PostCode,
		// 	City:       &addCustomerRequest.City,
		// 	Country:    &addCustomerRequest.Country,
		// },
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

	custparams := &stripe.CustomerParams{}
	custparams.DefaultSource = &pm.ID

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
	}

	c.JSON(http.StatusOK, gin.H{"complete": resp})
}

func getCustomer(c *gin.Context) {

}

func deleteCustomer(c *gin.Context) {

}

func getPrices(c *gin.Context) {
	stripe.Key = os.Getenv("STRIPE_KEY")
	var prods []*stripe.Product

	params := &stripe.ProductListParams{}
	i := product.List(params)
	for i.Next() {
		prod := i.Product()

		if prod.DefaultPrice != nil {
			priceparams := &stripe.PriceParams{}
			p, err := price.Get(prod.DefaultPrice.ID, priceparams)
			if err == nil {
				prod.DefaultPrice = p
			}
		}

		prods = append(prods, prod)
	}
	c.JSON(http.StatusOK, gin.H{"prices": prods})
}
