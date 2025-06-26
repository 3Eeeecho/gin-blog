package util

import (
	"errors"
	"time"

	"github.com/3Eeeecho/go-gin-example/pkg/setting"
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)
	claims := Claims{
		username,
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expireTime.Unix(),
			Issuer:    "gin-blog",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(getJWTSecret())
	return token, err
}

func ParseToken(tokenStr string) (*Claims, error) {
	if tokenStr == "" {
		return nil, errors.New("token cannot be empty")
	}

	if len(getJWTSecret()) == 0 {
		return nil, errors.New("JWT secret is not configured")
	}

	tokenClaims, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return getJWTSecret(), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

func getJWTSecret() []byte {
	return []byte(setting.AppSetting.JwtSecret)
}
