package auth

import (
	"database/sql"
	"errors"
	"github.com/labstack/echo/v4"
	httperror "github.com/mkolibaba/gophermart/internal/http/error"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type loginPayload struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Login(c echo.Context) error {
	var payload loginPayload
	if err := c.Bind(&payload); err != nil {
		// 400 — неверный формат запроса
		return httperror.InvalidRequestBody(err)
	}

	user, err := h.userService.UserGet(c.Request().Context(), payload.Login)

	var credentialsInvalid bool
	if err != nil {
		credentialsInvalid = errors.Is(err, sql.ErrNoRows)
	} else {
		credentialsInvalid = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)) != nil
	}

	if credentialsInvalid {
		// 401 — неверная пара логин/пароль
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}
	if err != nil {
		// 500 — внутренняя ошибка сервера
		return httperror.InternalServerError(err)
	}

	if err := h.setJWTCookie(c, user.Login); err != nil {
		// 500 — внутренняя ошибка сервера
		return httperror.InternalServerError(err)
	}

	// 200 — пользователь успешно аутентифицирован
	return c.NoContent(http.StatusOK)
}
