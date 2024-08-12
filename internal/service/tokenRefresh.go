package service

import (
	"Auth/auth"
	"Auth/internal/repository/postdb"
	"fmt"
)

func HandleRefreshTokenRequest(t auth.TokenPair) (*auth.TokenPair, error) {
	token := auth.TokenPair{}
	fmt.Println(t.AccessToken)
	guid, err := auth.ParseToken(t.AccessToken)
	if err != nil {
		return nil, err
	}
	oldRefresh, err := postdb.GetUser(guid)
	if err != nil {
		return nil, err
	}
	if err = auth.VerifyRefreshToken(oldRefresh, t.AccessToken); err != nil {
		return nil, err
	}
	a := auth.Auth{
		GUID: guid,
	}
	token = HandleTokenRequest(a)

	return &token, nil
}
