package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/mkolibaba/gophermart"
	"github.com/mkolibaba/gophermart/internal/auth"
	"net/http"
	"time"
)

const (
	jwtCookieName = "jwt"
	jwtPath       = "/"
)

type Handler struct {
	userService gophermart.UserService
}

func NewHandler(userService gophermart.UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (h *Handler) setJWTCookie(c echo.Context, userLogin string) error {
	token, err := auth.NewJWT(userLogin)
	if err != nil {
		return err
	}

	// TODO(improvement): выставлять expireCookie = jwt.ExpiresAt
	expireCookie := time.Now().Add(time.Hour)
	maxAge := int(expireCookie.Unix() - time.Now().Unix())
	c.SetCookie(&http.Cookie{
		Name:     jwtCookieName,
		Value:    token,
		MaxAge:   maxAge,
		Path:     jwtPath,
		HttpOnly: true,
	})

	return nil
}
