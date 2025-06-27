package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/mkolibaba/gophermart"
	"net/http"
	"time"
)

const (
	jwtCookieName = "jwt"
	jwtCookiePath = "/"
)

type Handler struct {
	JWTCookieName string
	JWTCookiePath string
	userService   gophermart.UserService
	authService   gophermart.AuthService
}

func NewHandler(userService gophermart.UserService, authService gophermart.AuthService) *Handler {
	return &Handler{
		JWTCookieName: jwtCookieName,
		JWTCookiePath: jwtCookiePath,
		userService:   userService,
		authService:   authService,
	}
}

func (h *Handler) setJWTCookie(c echo.Context, userLogin string) error {
	token, err := h.authService.NewJWT(userLogin)
	if err != nil {
		return err
	}

	// TODO(improvement): выставлять expireCookie = jwt.ExpiresAt
	expireCookie := time.Now().Add(time.Hour)
	maxAge := int(expireCookie.Unix() - time.Now().Unix())
	c.SetCookie(&http.Cookie{
		Name:     h.JWTCookieName,
		Value:    token,
		MaxAge:   maxAge,
		Path:     h.JWTCookiePath,
		HttpOnly: true,
	})

	return nil
}
