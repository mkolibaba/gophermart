package orders

import (
	"github.com/labstack/echo/v4"
	httperror "github.com/mkolibaba/gophermart/internal/http/error"
	"github.com/mkolibaba/gophermart/internal/http/security"
	"github.com/mkolibaba/gophermart/internal/slices"
	"github.com/mkolibaba/gophermart/postgres/gen"
	"net/http"
	"time"
)

type orderResponse struct {
	Number     string   `json:"number"`
	Status     string   `json:"status"`
	Accrual    *float64 `json:"accrual,omitempty"`
	UploadedAt string   `json:"uploaded_at"`
}

func (h *Handler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()

	userLogin := security.RetrieveUserLoginFromContext(c)

	result, err := h.orderService.OrderGetAll(ctx, userLogin)
	if err != nil {
		// 500 — внутренняя ошибка сервера
		return httperror.InternalServerError(err)
	}

	if len(result) == 0 {
		// 204 — нет данных для ответа
		return c.NoContent(http.StatusNoContent)
	}

	slices.SortByTimeDesc(result, func(order *postgres.Order) time.Time {
		return order.UploadedAt.Time
	})

	// 200 — успешная обработка запроса
	return c.JSON(http.StatusOK, slices.Transform(result, func(o *postgres.Order) orderResponse {
		return orderResponse{
			Number:     o.ID,
			UploadedAt: o.UploadedAt.Time.Format(time.RFC3339),
			Status:     string(o.AccrualStatus),
			Accrual:    o.AccrualPoints,
		}
	}))
}
