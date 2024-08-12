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
	Token    Auth      `json:"token"`
	ExpireAt time.Time `json:"expire_at"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

const (
	signingKey = "secret"
)

// Генерация JWT токена
func GenerateAccessToken(GUID Auth) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"exp":  time.Now().Add(time.Minute * 120).Unix(),
		"guid": GUID.GUID,
	})
	tokenString, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Генерация Refresh токена
func GenerateRefreshToken(t *RefreshToken) (string, error) {
	t.ExpireAt = time.Now().Add(time.Hour * 720)
	refreshToken := []byte(t.ExpireAt.String())
	hashRefreshToken, err := bcrypt.GenerateFromPassword(refreshToken, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	refresh := base64.StdEncoding.EncodeToString(hashRefreshToken)
	return refresh, nil
}

// Проверка Refresh токена
func VerifyRefreshToken(oldToken string, newToken string) error {
	hashRefreshToken, err := base64.StdEncoding.DecodeString(newToken)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(hashRefreshToken, []byte(oldToken))
	if err != nil {
		return err
	}
	return nil
}

func ParseToken(access string) (string, error) {
	fmt.Println(access)
	token, err := jwt.Parse(access, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что метод подписи соответствует ожидаемому
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Ожидался метод подписи HS512")
		}
		// Возвращаем секретный ключ для проверки подписи
		return []byte(signingKey), nil
	})
	if err != nil {
		panic(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Получаем значения из токена
		exp, _ := claims["exp"].(float64)
		guid, _ := claims["guid"].(string)

		fmt.Println("exp:", exp)
		fmt.Println("guid:", guid)
		return guid, nil
	} else {
		panic("Токен недействителен")
	}

	return "", nil
}
