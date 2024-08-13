package service

import (
	"Auth/auth"
	"Auth/internal/repository/postdb"
	"fmt"
	"log"
	"time"
)

// Выдача пары токенов
func HandleTokenRequest(a auth.Auth, clientIp string) (auth.TokenPair, error) {
	access, err := auth.GenerateAccessToken(a)
	if err != nil {
		log.Fatal(err)
	}
	dataRefreshToken := &auth.RefreshToken{
		Guid:     a,
		ExpireAt: time.Now().Add(time.Hour * 720).Unix(),
	}
	refreshToken, err := auth.GenerateRefreshToken(dataRefreshToken)
	if err != nil {
		log.Fatal(err)
	}
	tokenPair := auth.TokenPair{
		AccessToken:  access,
		RefreshToken: refreshToken,
	}
	res, err := postdb.GetUser(a.GUID)
	if err != nil {
		fmt.Println(err)
	}
	if res.RefreshToken == "" {
		postdb.AddUser(clientIp, dataRefreshToken, refreshToken)
	} else {
		log.Fatalf("User with guid %s already have refresh token", a.GUID)
	}
	return tokenPair, nil
}
