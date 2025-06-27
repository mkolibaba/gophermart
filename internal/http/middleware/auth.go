package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/mkolibaba/gophermart"
	httperror "github.com/mkolibaba/gophermart/internal/http/error"
	"strings"
)

func Auth(authService gophermart.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := extractJWTFromAuthorizationHeader(c)
			if tokenString == "" {
				tokenString = extractJWTFromCookie(c)
			}

			if tokenString == "" {
				// 401 — пользователь не авторизован
				return httperror.Unauthorized("no credentials provided")
			}

			claims, err := authService.GetClaims(tokenString)
			if err != nil {
				// 401 — пользователь не авторизован
				return httperror.Unauthorized("invalid credentials")
			}

			c.Set("claims", claims)
			return next(c)
		}
	}
}

func extractJWTFromAuthorizationHeader(c echo.Context) string {
	header := c.Request().Header.Get("Authorization")
	if header != "" {
		header = strings.TrimPrefix(header, "Bearer ")
	}
	return header
}

func extractJWTFromCookie(c echo.Context) string {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		return ""
	}
	return cookie.Value
}
