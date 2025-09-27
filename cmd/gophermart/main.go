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

func NewPostgresConnection(cfg *config.Config) (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), cfg.DatabaseURI)
}

func NewConfig() (*config.Config, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}
	stdlog.Printf("provided configuration: %+v", cfg)
	return cfg, nil
}

func RunMigrations(lc fx.Lifecycle, conn *pgx.Conn, logger *zap.SugaredLogger) {
	lc.Append(fx.StartHook(func(ctx context.Context) error {
		logger.Info("running database DDL migrations...")
		return migration.Run(ctx, conn)
	}))
}

func StartAccrualFetching(lc fx.Lifecycle, ordersService gophermart.OrderService) {
	lc.Append(fx.StartHook(ordersService.StartAccrualFetching))
}

func RunServer(*http.Server) {
}

func main() {
	fx.New(
		fx.Provide(
			NewConfig,
			NewSugaredLogger,
			NewPostgresConnection,
			fx.Annotate(postgres.NewDBX, fx.As(new(gophermart.Querier)), fx.As(new(gophermart.UserService))),
			fx.Annotate(accrual.NewClient, fx.As(new(gophermart.AccrualClient))),
			fx.Annotate(auth.NewService, fx.As(new(gophermart.AuthService))),
			fx.Annotate(withdraw.NewService, fx.As(new(gophermart.WithdrawService))),
			fx.Annotate(orders.NewService, fx.As(new(gophermart.OrderService))),
			http.NewServer,
		),
		fx.Invoke(
			RunMigrations,
			StartAccrualFetching,
			RunServer,
		),
	).Run()
}
