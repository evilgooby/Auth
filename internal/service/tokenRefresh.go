package service

import (
	"Auth/auth"
	"Auth/internal/repository/postdb"
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"time"
)

// Обновление токенов
func HandleRefreshTokenRequest(t auth.TokenPair, clientIp string) (auth.TokenPair, error) {
	token := auth.TokenPair{}
	GUID := auth.Auth{}
	guid, err := auth.ParseToken(t.AccessToken)
	if err != nil {
		return auth.TokenPair{}, err
	}
	GUID = auth.Auth{GUID: guid}
	bdRefresh, err := postdb.GetUser(guid)
	if err != nil {
		return auth.TokenPair{}, err
	}
	if err = VerifyExpiredToken(bdRefresh.ExpireAt); err != nil {
		return auth.TokenPair{}, fmt.Errorf("Token expired: %s", err)
	}
	if err = VerifyIP(clientIp, bdRefresh.Ip); err != nil {
		return auth.TokenPair{}, err
	}
	access, err := auth.GenerateAccessToken(GUID)
	if err != nil {
		return auth.TokenPair{}, fmt.Errorf("Failed to generate access token: %s", err)
	}
	dataRefreshToken := &auth.RefreshToken{
		Guid:     GUID,
		ExpireAt: time.Now().Add(time.Hour * 720).Unix(),
	}
	refreshToken, err := auth.GenerateRefreshToken(dataRefreshToken)
	if err != nil {
		return auth.TokenPair{}, err
	}
	token.RefreshToken = refreshToken
	token.AccessToken = access
	if err = auth.VerifyRefreshToken(token.RefreshToken, guid); err != nil {
		return auth.TokenPair{}, err
	}
	if err = auth.VerifyClientRefreshToken(bdRefresh.RefreshToken, t.RefreshToken); err != nil {
		return auth.TokenPair{}, err
	}
	if err = postdb.DeleteUser(guid); err != nil {
		return auth.TokenPair{}, err
	}
	if err = postdb.AddUser(clientIp, dataRefreshToken, refreshToken); err != nil {
		return auth.TokenPair{}, err
	}
	return token, nil
}

// Проверка не сменился ли IP адрес у клиента
func VerifyIP(ip string, ipBD string) error {
	if ip == ipBD {
		return nil
	} else {
		if err := godotenv.Load(); err != nil {
			log.Print("No .env file found")
		}
		myEmail := os.Getenv("MY_EMAIL")
		myPass := os.Getenv("MY_PASS")
		m := gomail.NewMessage()
		m.SetHeader("From", myEmail)
		m.SetHeader("To", "evilgooby1@gmail.com")
		m.SetHeader("Subject", "email warning")
		m.SetBody("text/plain", "email warning")

		d := gomail.NewDialer("smtp.gmail.com", 587, myEmail, myPass)

		if err := d.DialAndSend(m); err != nil {
			return fmt.Errorf("Failed to send email: %s", err)
		}

		fmt.Println("The letter was sent successfully")
		return fmt.Errorf("IP changed on client side")
	}
}

// Проверка не истек ли Refresh токен
func VerifyExpiredToken(token int64) error {
	if time.Now().Unix() < token {
		return nil
	} else {
		return fmt.Errorf("Token expired")
	}
}
