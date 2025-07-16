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

type Service struct {
	JWTSecretKey  string
	JWTTimeToLive time.Duration
	JWTIssuer     string
}

func NewService() *Service {
	return &Service{
		JWTIssuer:     issuer,
		JWTTimeToLive: ttl,
		JWTSecretKey:  secretKey,
	}
}

func (s *Service) NewJWT(login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   login,
		Issuer:    s.JWTIssuer,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.JWTTimeToLive)),
	})

	signed, err := token.SignedString([]byte(s.JWTSecretKey))
	if err != nil {
		return "", fmt.Errorf("auth: jwt signing: %w", err)
	}

	return signed, nil
}

func (s *Service) GetClaims(tokenString string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error: unexpected signing method")
		}
		return []byte(s.JWTSecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("auth: invalid token: %w", err)
	}

	if claims := token.Claims; token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("auth: invalid claims: %w", err)
}
