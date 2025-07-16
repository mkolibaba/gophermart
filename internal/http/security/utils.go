package security

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const claimsKey = "claims"

func StoreClaimsToContext(c echo.Context, claims jwt.Claims) {
	c.Set(claimsKey, claims)
}

func RetrieveClaimsFromContext(c echo.Context) jwt.Claims {
	return c.Get(claimsKey).(jwt.Claims)
}

func RetrieveUserLoginFromContext(c echo.Context) string {
	// имплементация всегда возвращает err = nil
	subject, _ := RetrieveClaimsFromContext(c).GetSubject()
	return subject
}
