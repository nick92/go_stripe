package stripe

import (
	"github.com/gin-gonic/gin"
)

func InitCustomerRoutes(router *gin.RouterGroup) {
	router.GET("/get", getCustomer)
	router.POST("/add", addCustomer)
	router.POST("/add_card", createNewCustomerPaymentMethod)
	router.DELETE("/delete", deleteCustomer)
	router.DELETE("/delete_card", removeCustomerPaymentMethod)
	router.GET("/payment_methods", getCustomerPaymentMethods)
	router.GET("/set_default", setCustomerDefaultPaymentMethods)
}

func InitSubscriptionRoutes(router *gin.RouterGroup) {
	router.DELETE("/cancel", canceCustomerSubscription)
}

func InitPriceRoutes(router *gin.RouterGroup) {
	router.GET("/get", getPrices)
}
