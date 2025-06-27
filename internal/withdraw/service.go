package withdraw

import (
	"context"
	"github.com/mkolibaba/gophermart"
	"github.com/mkolibaba/gophermart/postgres/gen"
)

type Service struct {
	userService gophermart.UserService
	querier     gophermart.Querier
}

func NewService(userService gophermart.UserService, querier gophermart.Querier) *Service {
	return &Service{
		userService: userService,
		querier:     querier,
	}
}

func (s *Service) WithdrawBalance(ctx context.Context, userLogin string, orderID string, sum float64) error {
	user, err := s.userService.UserGet(ctx, userLogin)
	if err != nil {
		return err
	}

	if user.AccrualBalance < sum {
		return gophermart.ErrInsufficientFunds
	}

	return s.querier.DoInTx(ctx, func(qtx postgres.Querier) error {
		// сохраняем списание
		if err := qtx.WithdrawalSave(ctx, postgres.WithdrawalSaveParams{
			OrderNumber: orderID,
			UserLogin:   userLogin,
			Sum:         sum,
		}); err != nil {
			return err
		}

		// обновляем баланс пользователя
		return qtx.UserAddAccrualBalance(ctx, -sum, userLogin)
	})
}
