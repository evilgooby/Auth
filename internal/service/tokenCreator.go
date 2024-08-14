package service

import (
	"Auth/internal/auth"
	"Auth/internal/middleware"
	"Auth/internal/repository/postdb"
	"github.com/gin-gonic/gin"
	"time"
)

// Выдача пары токенов
func HandleTokenRequest(a auth.Auth, clientIp string, c *gin.Context) (auth.TokenPair, error) {
	access, err := auth.GenerateAccessToken(a)
	if err != nil {
		return auth.TokenPair{}, c.Error(middleware.ErrInternalServerError)
	}
	dataRefreshToken := &auth.RefreshToken{
		Guid:     a,
		ExpireAt: time.Now().Add(time.Hour * 720).Unix(),
	}
	refreshToken, err := auth.GenerateRefreshToken(dataRefreshToken)
	if err != nil {
		return auth.TokenPair{}, c.Error(middleware.ErrInternalServerError)
	}
	tokenPair := auth.TokenPair{
		AccessToken:  access,
		RefreshToken: refreshToken,
	}
	res, err := postdb.GetUser(a.GUID)
	if err != nil {
		return auth.TokenPair{}, c.Error(err)
	}
	if res.RefreshToken == "" {
		postdb.AddUser(clientIp, dataRefreshToken, refreshToken)
	} else {
		return auth.TokenPair{}, c.Error(err)
	}
	return tokenPair, nil
}
