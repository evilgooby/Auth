package service

import (
	"Auth/auth"
	"Auth/internal/repository/postdb"
	"fmt"
	"time"
)

// Выдача пары токенов
func HandleTokenRequest(a auth.Auth, clientIp string) (auth.TokenPair, error) {
	access, err := auth.GenerateAccessToken(a)
	if err != nil {
		return auth.TokenPair{}, err
	}
	dataRefreshToken := &auth.RefreshToken{
		Guid:     a,
		ExpireAt: time.Now().Add(time.Hour * 720).Unix(),
	}
	refreshToken, err := auth.GenerateRefreshToken(dataRefreshToken)
	if err != nil {
		return auth.TokenPair{}, err
	}
	tokenPair := auth.TokenPair{
		AccessToken:  access,
		RefreshToken: refreshToken,
	}
	res, err := postdb.GetUser(a.GUID)
	if err != nil {
		return auth.TokenPair{}, err
	}
	if res.RefreshToken == "" {
		postdb.AddUser(clientIp, dataRefreshToken, refreshToken)
	} else {
		return auth.TokenPair{}, fmt.Errorf("User with guid %s already have refresh token", a.GUID)
	}
	return tokenPair, nil
}
