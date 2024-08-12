package controller

import (
	"Auth/auth"
	"Auth/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func AddToken(c *gin.Context) {
	var aut auth.Auth
	if err := c.ShouldBindJSON(&aut); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON provided"})
		return
	}
	log.Println(aut)
	token := service.HandleTokenRequest(aut)
	c.JSON(http.StatusOK, gin.H{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})
}

func RefreshToken(c *gin.Context) {
	var tokenPair auth.TokenPair
	if err := c.ShouldBindJSON(&tokenPair); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON provided"})
		return
	}
	fmt.Println(tokenPair)
	token, err := service.HandleRefreshTokenRequest(tokenPair)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(token)
	c.JSON(http.StatusOK, gin.H{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})
}
