package main

import (
	"os"
	"trucknav/stripe"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		// common.NewBusinessError(common.System, "PORT variable must be set", 0)
		return
	}

	if os.Getenv("MODE") == "Release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := SetupRouter()
	router.Run(":" + port)
}

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
