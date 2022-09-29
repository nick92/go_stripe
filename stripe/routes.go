package stripe

import (
	"github.com/gin-gonic/gin"
)

func InitCustomerRoutes(router *gin.RouterGroup) {
	router.GET("/get", getCustomer)
	router.POST("/add", addCustomer)
	router.DELETE("/delete", deleteCustomer)
	router.GET("/payment_methods", getCustomerPaymentMethods)
}

func InitPriceRoutes(router *gin.RouterGroup) {
	router.GET("/get", getPrices)
	router.POST("/add", createPaymentMethod)
}
