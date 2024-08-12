package service

import (
	"Auth/auth"
	"Auth/internal/repository/postdb"
	"fmt"
	"log"
	"time"
)

func HandleRefreshTokenRequest(t auth.TokenPair, clientIp string) (*auth.TokenPair, error) {
	token := auth.TokenPair{}
	GUID := auth.Auth{}
	guid, err := auth.ParseToken(t.AccessToken)
	if err != nil {
		return nil, err
	}
	GUID = auth.Auth{GUID: guid}
	bdRefresh, err := postdb.GetUser(guid)
	if err != nil {
		return nil, err
	}
	err = VerifyIP(clientIp, bdRefresh.Ip)
	if err != nil {
		return nil, err
	}
	access, err := auth.GenerateAccessToken(GUID)
	if err != nil {
		panic(err)
	}
	dataRefreshToken := &auth.RefreshToken{
		Guid:     GUID,
		ExpireAt: time.Now().Add(time.Hour * 720).Unix(),
	}
	refreshToken, err := auth.GenerateRefreshToken(dataRefreshToken)
	if err != nil {
		log.Fatal(err)
	}
	token.RefreshToken = refreshToken
	token.AccessToken = access
	err = auth.VerifyRefreshToken(token.RefreshToken, guid)
	if err != nil {
		return nil, err
	}
	err = auth.VerifyClientRefreshToken(bdRefresh.RefreshToken, t.RefreshToken)
	if err != nil {
		return nil, err
	}
	err = postdb.DeleteUser(guid)
	if err != nil {
		return nil, err
	}
	err = postdb.AddUser(clientIp, dataRefreshToken, refreshToken)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func VerifyIP(ip string, ipBD string) error {
	if ip == ipBD {
		return nil
	} else {
		return fmt.Errorf("IP изменен на стороне клиента")
	}

}
