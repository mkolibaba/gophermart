package gophermart

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mkolibaba/gophermart/postgres/gen"
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
)

type (
	AuthService interface {
		// NewJWT создает JWT для переданного логина
		NewJWT(login string) (string, error)

		// GetClaims получает JWT Claims из переданного токена
		GetClaims(tokenString string) (jwt.Claims, error)
	}

	OrderService interface {
		// OrderGet возвращает заказ по переданному id
		OrderGet(ctx context.Context, id string) (*postgres.Order, error)

		// OrderGetAll возвращает все заказы пользователя
		OrderGetAll(ctx context.Context, userLogin string) ([]*postgres.Order, error)

		// OrderSave создает заказ с переданными параметрами
		OrderSave(ctx context.Context, id string, userLogin string) error
	}

	UserService interface {
		// UserGet возвращает пользователя по переданному login
		UserGet(ctx context.Context, login string) (*postgres.User, error)

		// UserSave создает пользователя с переданными параметрами
		UserSave(ctx context.Context, login string, password string) error

		// UserAddAccrualBalance начисляет на баланс пользователя заданную сумму
		UserAddAccrualBalance(ctx context.Context, sum float64, login string) error
	}

	BalanceService interface {
		// BalanceGet возвращает информацию о текущем балансе и сумме списаний пользователя
		BalanceGet(ctx context.Context, userLogin string) (*postgres.BalanceGetRow, error)

		// WithdrawalGetAll возвращает все списания пользователя
		WithdrawalGetAll(ctx context.Context, userLogin string) ([]*postgres.Withdrawal, error)
	}

	WithdrawService interface {
		// WithdrawBalance списывает бонусы со счета клиента, сохраняя в истории данное списание. Возвращает
		// ErrInsufficientFunds в случае нехватки бонусов
		WithdrawBalance(ctx context.Context, userLogin string, orderID string, sum float64) error
	}

	Querier interface {
		postgres.Querier

		// DoInTx выполняет переданную функцию в транзакции qtx
		DoInTx(ctx context.Context, fn func(qtx postgres.Querier) error) error
	}

	AccrualOrder struct {
		Order   string   `json:"order"`
		Status  string   `json:"status"`
		Accrual *float64 `json:"accrual,omitempty"`
	}

	AccrualClient interface {
		// GetOrder получает данные по заказу из системы расчета баллов лояльности
		GetOrder(ctx context.Context, number string) (*AccrualOrder, error)
	}
)
