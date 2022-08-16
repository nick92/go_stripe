package main

import (
	"os"
	"trucknav/router"

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

	router := router.SetupRouter()
	router.Run(":" + port)
}
