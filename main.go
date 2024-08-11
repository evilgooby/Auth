package main

import (
	"Authentication/auth"
	"Authentication/postdb"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"net/http"
	"time"
)

func handleTokenRequest(a *auth.Auth) auth.TokenPair {
	fmt.Println(a.GUID)
	acess, err := auth.GenerateAccessToken(a)
	if err != nil {
		panic(err)
	}
	refreshToken := &auth.RefreshToken{
		Token:    *a,
		ExpireAt: time.Now().Add(time.Hour * 720),
	}
	refresh, err := auth.GenerateRefreshToken(refreshToken)
	if err != nil {
		panic(err)
	}
	tokenPair := auth.TokenPair{
		AccessToken:  acess,
		RefreshToken: refresh,
	}
	res, err := postdb.GetUser(a.GUID)
	if err != nil {
		fmt.Println(err)
	}
	if res == "" {
		postdb.AddUser(a.GUID, refresh)
	} else {
		postdb.DeleteUser(a.GUID)
		postdb.AddUser(a.GUID, refresh)
	}
	return tokenPair
}

// Обработчик запроса на обновление токена
func handleRefreshTokenRequest(t *auth.TokenPair) auth.TokenPair {
	tok := auth.TokenPair{}
	fmt.Println(t.AccessToken)
	guid, err := auth.ParseToken(t.AccessToken)
	if err != nil {
		panic(err)
	}
	oldRefresh, err := postdb.GetUser(guid)
	if err != nil {
		panic(err)
	}
	exam, err := auth.VerifyRefreshToken(oldRefresh, t.AccessToken)
	if err != nil {
		panic(err)
	}
	if exam {
		a := &auth.Auth{
			GUID: guid,
		}
		tok = handleTokenRequest(a)

	} else {
		fmt.Errorf("Токен изменен")
	}
	return tok
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to authentication API",
		})
	})

	r.POST("/addToken", func(c *gin.Context) {
		var aut auth.Auth
		if err := c.ShouldBindQuery(&aut); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON provided"})
			return
		}
		fmt.Println(&aut.GUID)
		token := handleTokenRequest(&aut)
		c.JSON(http.StatusOK, gin.H{
			"access_token":  token.AccessToken,
			"refresh_token": token.RefreshToken,
		})
	})

	r.POST("/refreshToken", func(c *gin.Context) {
		var tokenPair auth.TokenPair
		if err := c.ShouldBindQuery(&tokenPair); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON provided"})
			return
		}
		fmt.Println(tokenPair)
		token := handleRefreshTokenRequest(&tokenPair)
		fmt.Println(token)
		c.JSON(http.StatusOK, gin.H{
			"access_token":  token.AccessToken,
			"refresh_token": token.RefreshToken,
		})
	})

	err := r.Run(":8080")
	if err != nil {
		panic(fmt.Sprintf("Failed to start the web server - Error: %v", err))
	}

}
