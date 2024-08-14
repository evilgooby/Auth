package controller

import (
	"Auth/internal/middleware"
	"github.com/gin-gonic/gin"
	"log"
)

const port = ":8080"

func Registry() {
	r := gin.Default()
	r.Use(middleware.ErrorHandler)

	r.POST("/addToken", AddToken)
	r.POST("/refreshToken", RefreshToken)

	err := r.Run(port)
	if err != nil {
		log.Fatalf("Failed to start the web server - Error: %v", err)
	}
}
