package stripe

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/product"
)

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
