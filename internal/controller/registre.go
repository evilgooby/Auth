package controller

import (
	"github.com/gin-gonic/gin"
	"log"
)

func Registre() {
	r := gin.Default()

	r.POST("/addToken", AddToken)
	r.POST("/refreshToken", RefreshToken)

	err := r.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start the web server - Error: %v", err)
	}
}
