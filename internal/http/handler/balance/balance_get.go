package balance

import (
	"github.com/labstack/echo/v4"
	httperror "github.com/mkolibaba/gophermart/internal/http/error"
	"github.com/mkolibaba/gophermart/internal/http/security"
	"net/http"
)

func (h *Handler) Get(c echo.Context) error {
	userLogin := security.RetrieveUserLoginFromContext(c)

	result, err := h.balanceService.BalanceGet(c.Request().Context(), userLogin)
	if err != nil {
		// 500 — внутренняя ошибка сервера
		return httperror.InternalServerError(err)
	}

	// 200 — успешная обработка запроса
	return c.JSON(http.StatusOK, echo.Map{
		"current":   result.Current,
		"withdrawn": result.Withdrawn,
	})
}
