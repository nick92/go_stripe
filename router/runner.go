package router

import (
	"trucknav/stripe"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// create new router for site routes and memory store
	router := gin.Default()

	// Setup device routes
	stripeapi := router.Group("/api/customer")
	stripe.InitCustomerRoutes(stripeapi)

	pricesapi := router.Group("/api/prices")
	stripe.InitPriceRoutes(pricesapi)

	return router
}
