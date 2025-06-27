package auth

import (
	"github.com/labstack/echo/v4"
	httperror "github.com/mkolibaba/gophermart/internal/http/error"
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

	exists, err := h.userService.UserExists(ctx, payload.Login)
	if err != nil {
		// 500 — внутренняя ошибка сервера
		return httperror.InternalServerError(err)
	}
	if exists {
		// 409 — логин уже занят
		// TODO(improvement): кейс можно совместить с userService.UserSave и проверять возвращаемый err
		return echo.NewHTTPError(http.StatusConflict, "user already exists")
	}

	// TODO(improvement): не хранить пароли в открытом виде
	if err := h.userService.UserSave(ctx, payload.Login, payload.Password); err != nil {
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
