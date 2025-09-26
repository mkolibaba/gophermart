package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/mkolibaba/gophermart"
	"github.com/mkolibaba/gophermart/internal/auth"
	"github.com/mkolibaba/gophermart/internal/config"
	"github.com/mkolibaba/gophermart/internal/http"
	"github.com/mkolibaba/gophermart/internal/http/client/accrual"
	"github.com/mkolibaba/gophermart/internal/orders"
	"github.com/mkolibaba/gophermart/internal/withdraw"
	"github.com/mkolibaba/gophermart/postgres"
	"github.com/mkolibaba/gophermart/postgres/migration"
	"go.uber.org/fx"
	"go.uber.org/zap"
	stdlog "log"
)

func NewSugaredLogger() *zap.SugaredLogger {
	unsugaredLogger, err := zap.NewDevelopment()
	if err != nil {
		stdlog.Fatal(err)
	}
	return unsugaredLogger.Sugar()
}

func NewPostgresConnection(ctx context.Context, cfg *config.Config) (*pgx.Conn, error) {
	return pgx.Connect(ctx, cfg.DatabaseURI)
}

func NewHTTPService(
	lc fx.Lifecycle,
	cfg *config.Config,
	logger *zap.SugaredLogger,
	dbx gophermart.Querier,
	authService gophermart.AuthService,
	ordersService gophermart.OrderService,
	withdrawService gophermart.WithdrawService,
) *http.Server {
	h := &http.Server{
		Address:         cfg.RunAddress,
		Logger:          logger,
		Querier:         dbx,
		AuthService:     authService,
		OrderService:    ordersService,
		WithdrawService: withdrawService,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go h.Start(ctx)
			return nil
		},
	})
	return h
}

func main() {
	fx.New(
		fx.Provide(config.New),
		fx.Provide(NewSugaredLogger),
		fx.Provide(NewPostgresConnection),
		fx.Provide(postgres.NewDBX),
		fx.Provide(accrual.NewClient),
		fx.Provide(auth.NewService),
		fx.Provide(withdraw.NewService),
		fx.Provide(orders.NewService),
		fx.Invoke(func(ctx context.Context, conn *pgx.Conn, logger *zap.SugaredLogger) {
			logger.Info("running database DDL migrations...")
			if err := runMigrations(ctx, conn); err != nil {
				logger.Fatalf("failed to run database ddl migrations: %s", err)
			}
		}),
		fx.Invoke(func(ctx context.Context, ordersService *orders.Service) {
			ordersService.StartAccrualFetching(ctx)
		}),
	).Run()
}

func runMigrations(ctx context.Context, conn *pgx.Conn) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := conn.Exec(ctx, migration.DDL); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
