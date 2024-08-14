package service

import (
	"Auth/internal/auth"
	"Auth/internal/middleware"
	"Auth/internal/repository/postdb"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"os"
	"time"
)

// Обновление токенов
func HandleRefreshTokenRequest(c *gin.Context, t auth.TokenPair, clientIp string) (*auth.TokenPair, error) {
	token := auth.TokenPair{}
	guid, err := auth.ParseToken(t.AccessToken)
	if err != nil {
		return nil, c.Error(err)
	}
	bdRefresh, err := postdb.GetUser(guid.GUID)
	if err != nil {
		return nil, c.Error(err)
	}
	if err = VerifyExpiredToken(bdRefresh.ExpireAt); err != nil {
		return nil, c.Error(err)
	}
	if err = VerifyIP(clientIp, bdRefresh.Ip); err != nil {
		return nil, c.Error(err)
	}
	access, err := auth.GenerateAccessToken(*guid)
	if err != nil {
		return nil, c.Error(err)
	}
	dataRefreshToken := &auth.RefreshToken{
		Guid:     *guid,
		ExpireAt: time.Now().Add(time.Hour * 720).Unix(),
	}
	refreshToken, err := auth.GenerateRefreshToken(dataRefreshToken)
	if err != nil {
		return nil, c.Error(err)
	}
	token.RefreshToken = refreshToken
	token.AccessToken = access
	if err = auth.VerifyRefreshToken(token.RefreshToken, guid.GUID); err != nil {
		return nil, c.Error(err)
	}
	if err = auth.VerifyClientRefreshToken(bdRefresh.RefreshToken, t.RefreshToken); err != nil {
		return nil, c.Error(err)
	}
	if err = postdb.UpdateUser(clientIp, dataRefreshToken, refreshToken); err != nil {
		return nil, c.Error(err)
	}
	return &token, nil
}

// Проверка не сменился IP адрес у клиента
func VerifyIP(ip string, ipBD string) error {
	if ip == ipBD {
		return nil
	} else {
		myEmail := os.Getenv("MY_EMAIL")
		myPass := os.Getenv("MY_PASS")
		m := gomail.NewMessage()
		m.SetHeader("From", myEmail)
		m.SetHeader("To", "evilgooby1@gmail.com")
		m.SetHeader("Subject", "email warning")
		m.SetBody("text/plain", "email warning")

		d := gomail.NewDialer("smtp.gmail.com", 587, myEmail, myPass)

		if err := d.DialAndSend(m); err != nil {
			return middleware.ErrEmailSend
		}

		fmt.Println("The letter was sent successfully")
		return middleware.ErrIpChange
	}
}

// Проверка не истек ли Refresh токен
func VerifyExpiredToken(token int64) error {
	if time.Now().Unix() < token {
		return nil
	} else {
		return middleware.ErrExpiredToken
	}
}
