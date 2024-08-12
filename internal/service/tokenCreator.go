package service

import (
	"Auth/auth"
	"Auth/internal/repository/postdb"
	"fmt"
	"log"
	"time"
)

func HandleTokenRequest(a auth.Auth, clientIp string) auth.TokenPair {
	acess, err := auth.GenerateAccessToken(a)
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
		AccessToken:  acess,
		RefreshToken: refreshToken,
	}
	res, err := postdb.GetUser(a.GUID)
	if err != nil {
		fmt.Println(err)
	}
	if res.RefreshToken == "" {
		postdb.AddUser(clientIp, dataRefreshToken, refreshToken)
	} else {
		postdb.DeleteUser(a.GUID)
		postdb.AddUser(clientIp, dataRefreshToken, refreshToken)
	}
	return tokenPair
}
