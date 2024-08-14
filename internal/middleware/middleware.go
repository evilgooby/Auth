package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	ErrNotFound            = fmt.Errorf("not found")
	ErrInternalServerError = fmt.Errorf("internal server error")
	ErrHaveRefreshToken    = fmt.Errorf("user with guid already have refresh token")
	ErrInvalidToken        = fmt.Errorf("invalid token")
	ErrExpiredToken        = fmt.Errorf("token expired")
	ErrDB                  = fmt.Errorf("database error")
)

func ErrorHandler(c *gin.Context) {
	c.Next()

	for _, err := range c.Errors {
		switch err.Err {
		case ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": ErrNotFound.Error()})
			return
		case ErrInternalServerError:
			c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInternalServerError.Error()})
			return
		case ErrHaveRefreshToken:
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrHaveRefreshToken.Error()})
			return
		case ErrInvalidToken:
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidToken.Error()})
			return
		case ErrExpiredToken:
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrExpiredToken.Error()})
			return
		case ErrDB:
			c.JSON(http.StatusInternalServerError, gin.H{"error": ErrDB.Error()})
			return
		}
	}

	c.JSON(http.StatusInternalServerError, "")
}
