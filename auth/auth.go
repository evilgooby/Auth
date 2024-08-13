package auth

import (
	"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Auth struct {
	GUID string `json:"guid"`
}
type RefreshToken struct {
	Guid     Auth  `json:"token"`
	ExpireAt int64 `json:"expire_at"`
}
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type ClientRefreshToken struct {
	Ip           string `json:"ip"`
	RefreshToken string `json:"refresh_token"`
	ExpireAt     int64  `json:"expires_in"`
}

const (
	signingKey = "secret"
)

// Генерация JWT токена
func GenerateAccessToken(GUID Auth) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"ExpireAt": time.Now().Add(time.Minute * 120).Unix(),
		"Guid":     GUID.GUID,
	})
	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Генерация Refresh токена
func GenerateRefreshToken(t *RefreshToken) (string, error) {
	refreshToken := []byte(t.Guid.GUID)
	hashRefreshToken, err := bcrypt.GenerateFromPassword(refreshToken, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	refresh := base64.StdEncoding.EncodeToString(hashRefreshToken)
	return refresh, nil
}

// Проверка Refresh токена
func VerifyRefreshToken(newToken string, guid string) error {
	hashRefreshToken, err := base64.StdEncoding.DecodeString(newToken)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(hashRefreshToken, []byte(guid))
	if err != nil {
		return err
	}
	return nil
}

// Проверка Refresh токена на изменения на стороне клиента
func VerifyClientRefreshToken(tokenBD, tokenClient string) error {
	if tokenClient == tokenBD {
		return nil
	} else {
		return fmt.Errorf("Token modified on client side")
	}
}

// Парсинг Access токена для получения GUID
func ParseToken(access string) (string, error) {
	token, err := jwt.Parse(access, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что метод подписи соответствует ожидаемому
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("HS512 update method expected")
		}
		// Возвращаем секретный ключ для проверки подписи
		return []byte(signingKey), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		guid, _ := claims["Guid"].(string)
		return guid, nil
	} else {
		return "", fmt.Errorf("Token is invalid")
	}
}
