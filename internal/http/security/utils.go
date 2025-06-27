package security

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func RetrieveUserLoginFromContext(c echo.Context) string {
	claims := c.Get("claims")
	// имплементация всегда возвращает err = nil
	subject, _ := claims.(jwt.RegisteredClaims).GetSubject()
	return subject
}
