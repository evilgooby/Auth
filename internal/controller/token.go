package controller

import (
	"Auth/internal/auth"
	"Auth/internal/middleware"
	"Auth/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddToken(c *gin.Context) {
	var aut auth.Auth
	if err := c.ShouldBindJSON(&aut); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON provided"})
		return
	}
	clientIP := c.ClientIP()
	token, err := service.HandleTokenRequest(c, aut, clientIP)
	if err != nil {
		middleware.ErrorHandler(c)
		return
	}
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
	clientIP := c.ClientIP()
	token, err := service.HandleRefreshTokenRequest(c, tokenPair, clientIP)
	if err != nil {
		middleware.ErrorHandler(c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"access_token":  token.AccessToken,
		"refresh_token": token.RefreshToken,
	})
}
