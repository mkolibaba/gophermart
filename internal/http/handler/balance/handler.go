package balance

import "github.com/mkolibaba/gophermart"

type Handler struct {
	balanceService    gophermart.BalanceService
	withdrawalService gophermart.WithdrawService
}

func NewHandler(balanceService gophermart.BalanceService, withdrawalService gophermart.WithdrawService) *Handler {
	return &Handler{
		balanceService:    balanceService,
		withdrawalService: withdrawalService,
	}
}
