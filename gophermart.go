package gophermart

import (
	"context"
	"errors"
	"github.com/mkolibaba/gophermart/postgres/gen"
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
)

type (
	OrderService interface {
		OrderGet(ctx context.Context, id string) (*postgres.Order, error)
		OrderGetAll(ctx context.Context, userLogin string) ([]*postgres.Order, error)
		OrderSave(ctx context.Context, id string, userLogin string) error
	}

	UserService interface {
		UserGet(ctx context.Context, login string) (*postgres.User, error)
		UserSave(ctx context.Context, login string, password string) error
		UserExists(ctx context.Context, login string) (bool, error)
		UserGetForLoginAndPassword(ctx context.Context, login, password string) (*postgres.User, error)

		// UserAddAccrualBalance начисляет на баланс пользователя заданную сумму
		UserAddAccrualBalance(ctx context.Context, sum float64, login string) error
	}

	BalanceService interface {
		BalanceGet(ctx context.Context, userLogin string) (*postgres.BalanceGetRow, error)
		WithdrawalGetAll(ctx context.Context, userLogin string) ([]*postgres.Withdrawal, error)
	}

	WithdrawService interface {
		// WithdrawBalance списывает бонусы со счета клиента, сохраняя в истории данное списание. Возвращает
		// ErrInsufficientFunds в случае нехватки бонусов
		WithdrawBalance(ctx context.Context, userLogin string, orderID string, sum float64) error
	}

	Querier interface {
		postgres.Querier
		DoInTx(ctx context.Context, fn func(qtx postgres.Querier) error) error
	}

	AccrualOrder struct {
		Order   string   `json:"order"`
		Status  string   `json:"status"`
		Accrual *float64 `json:"accrual,omitempty"`
	}

	AccrualClient interface {
		GetOrder(ctx context.Context, number string) (*AccrualOrder, error)
	}
)
