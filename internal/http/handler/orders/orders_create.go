package orders

import (
	"database/sql"
	"errors"
	"github.com/labstack/echo/v4"
	httperror "github.com/mkolibaba/gophermart/internal/http/error"
	"github.com/mkolibaba/gophermart/internal/http/security"
	"github.com/mkolibaba/gophermart/internal/validation"
	"io"
	"net/http"
)

func (h *Handler) Create(c echo.Context) error {
	payload, err := io.ReadAll(c.Request().Body)
	if err != nil {
		// 400 — неверный формат запроса
		return httperror.InvalidRequestBody(err)
	}
	orderId := string(payload)

	if !validation.Luhn(orderId) {
		// 422 — неверный формат номера заказа
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "invalid order id")
	}

	userLogin := security.RetrieveUserLoginFromContext(c)

	order, err := h.orderService.OrderGet(c.Request().Context(), orderId)
	if err == nil {
		if order.UserLogin != userLogin {
			// 409 — номер заказа уже был загружен другим пользователем
			return echo.NewHTTPError(http.StatusConflict, "order has been created by another user")
		}

		// 200 — номер заказа уже был загружен этим пользователем
		return c.NoContent(http.StatusOK)
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return httperror.InternalServerError(err)
	}

	if err := h.orderService.OrderSave(c.Request().Context(), orderId, userLogin); err != nil {
		return httperror.InternalServerError(err)
	}

	// 202 — новый номер заказа принят в обработку
	return c.NoContent(http.StatusAccepted)
}
