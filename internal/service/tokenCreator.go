package service

import (
	"Auth/auth"
	"Auth/internal/repository/postdb"
	"fmt"
	"time"
)

func HandleTokenRequest(a auth.Auth) auth.TokenPair {
	fmt.Println(a.GUID)
	acess, err := auth.GenerateAccessToken(a)
	if err != nil {
		panic(err)
	}
	refreshToken := &auth.RefreshToken{
		Token:    a,
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
