package balance

import (
	"github.com/labstack/echo/v4"
	httperror "github.com/mkolibaba/gophermart/internal/http/error"
	"github.com/mkolibaba/gophermart/internal/http/security"
	"github.com/mkolibaba/gophermart/internal/slices"
	"github.com/mkolibaba/gophermart/postgres/gen"
	"net/http"
	"time"
)

type withdrawalResponse struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

func (h *Handler) GetAllWithdrawals(c echo.Context) error {
	userLogin := security.RetrieveUserLoginFromContext(c)

	withdrawals, err := h.balanceService.WithdrawalGetAll(c.Request().Context(), userLogin)
	if err != nil {
		// 500 — внутренняя ошибка сервера
		return httperror.InternalServerError(err)
	}

	if len(withdrawals) == 0 {
		// 204 — нет ни одного списания
		return c.NoContent(http.StatusNoContent)
	}

	slices.SortByTimeDesc(withdrawals, func(w *postgres.Withdrawal) time.Time {
		return w.ProcessedAt.Time
	})

	// 200 — успешная обработка запроса
	return c.JSON(http.StatusOK, slices.Transform(withdrawals, func(s *postgres.Withdrawal) withdrawalResponse {
		return withdrawalResponse{
			Order:       s.OrderNumber,
			Sum:         s.Sum,
			ProcessedAt: s.ProcessedAt.Time.Format(time.RFC3339),
		}
	}))
}
