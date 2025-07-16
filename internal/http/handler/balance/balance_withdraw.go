package balance

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/mkolibaba/gophermart"
	httperror "github.com/mkolibaba/gophermart/internal/http/error"
	"github.com/mkolibaba/gophermart/internal/http/security"
	"github.com/mkolibaba/gophermart/internal/validation"
	"net/http"
)

type balanceWithdrawPayload struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func (h *Handler) WithdrawBalance(c echo.Context) error {
	ctx := c.Request().Context()

	userLogin := security.RetrieveUserLoginFromContext(c)

	var payload balanceWithdrawPayload
	if err := c.Bind(&payload); err != nil {
		// 400 — неверный формат запроса
		return httperror.InvalidRequestBody(err)
	}

	if !validation.Luhn(payload.Order) {
		// 422 — неверный формат номера заказа
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "invalid order id")
	}

	if err := h.withdrawalService.WithdrawBalance(ctx, userLogin, payload.Order, payload.Sum); err != nil {
		if errors.Is(err, gophermart.ErrInsufficientFunds) {
			// 402 — на счету недостаточно средств
			return echo.NewHTTPError(http.StatusPaymentRequired)
		}

		// 500 — внутренняя ошибка сервера
		return httperror.InternalServerError(err)
	}

	// 200 — успешная обработка запроса
	return c.NoContent(http.StatusOK)
}
