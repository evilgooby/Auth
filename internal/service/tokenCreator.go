package service

import (
	"Auth/internal/auth"
	"Auth/internal/repository/postdb"
	"github.com/gin-gonic/gin"
	"time"
)

// Выдача пары токенов
func HandleTokenRequest(c *gin.Context, a auth.Auth, clientIp string) (*auth.TokenPair, error) {
	access, err := auth.GenerateAccessToken(a)
	if err != nil {
		return nil, c.Error(err)
	}
	dataRefreshToken := &auth.RefreshToken{
		Guid:     a,
		ExpireAt: time.Now().Add(time.Hour * 720).Unix(),
	}
	refreshToken, err := auth.GenerateRefreshToken(dataRefreshToken)
	if err != nil {
		return nil, c.Error(err)
	}
	tokenPair := auth.TokenPair{
		AccessToken:  access,
		RefreshToken: refreshToken,
	}
	res, err := postdb.VerifyUser(a.GUID)
	if err != nil {
		return nil, c.Error(err)
	}
	if res.RefreshToken == "" {
		err = postdb.AddUser(clientIp, dataRefreshToken, refreshToken)
		if err != nil {
			return nil, c.Error(err)
		}
	}
	return &tokenPair, nil
}
