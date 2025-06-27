package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	secretKey = "GOPHERMART_APP"
	ttl       = 1 * time.Hour
	issuer    = "gophermart"
)

func NewJWT(login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   login,
		Issuer:    issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
	})

	signed, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func GetClaims(tokenString string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error: unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("auth: invalid token: %w", err)
	}

	if claims := token.Claims; token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("auth: invalid claims: %w", err)
}
