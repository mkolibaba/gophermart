package auth

import (
	"github.com/labstack/echo/v4"
	httperror "github.com/mkolibaba/gophermart/internal/http/error"
	"github.com/mkolibaba/gophermart/postgres"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type registerPayload struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Register(c echo.Context) error {
	ctx := c.Request().Context()

	var payload registerPayload
	if err := c.Bind(&payload); err != nil {
		// 400 — неверный формат запроса
		return httperror.InvalidRequestBody(err)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		// 500 — внутренняя ошибка сервера
		return httperror.InternalServerError(err)
	}

	if err := h.userService.UserSave(ctx, payload.Login, string(password)); err != nil {
		if postgres.IsUniqueViolationError(err) {
			// 409 — логин уже занят
			return echo.NewHTTPError(http.StatusConflict, "user already exists")
		}

		// 500 — внутренняя ошибка сервера
		return httperror.InternalServerError(err)
	}

	if err := h.setJWTCookie(c, payload.Login); err != nil {
		// 500 — внутренняя ошибка сервера
		return httperror.InternalServerError(err)
	}

	// 200 — пользователь успешно зарегистрирован и аутентифицирован
	return c.NoContent(http.StatusOK)
}
