package controller

import (
	"Auth/internal/middleware"
	"Auth/internal/repository/postdb"
	"github.com/gin-gonic/gin"
	"log"
)

func Registre() {
	r := gin.New()
	r.Use(middleware.ErrorHandler)

	r.POST("/addToken", AddToken)
	r.POST("/refreshToken", RefreshToken)
	err := postdb.InitialPostgres()
	if err != nil {
		log.Println("Failed to start the web server - Error: %v", err)
	}
	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start the web server - Error: %v", err)
	}
}
