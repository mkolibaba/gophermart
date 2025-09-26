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

func NewPostgresConnection(lc fx.Lifecycle, cfg *config.Config) *pgx.Conn {
	var conn *pgx.Conn
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			conn, err = pgx.Connect(ctx, cfg.DatabaseURI)
			return err
		},
	})
	return conn
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
		fx.Provide(
			config.New,
			NewSugaredLogger,
			NewPostgresConnection,
			fx.Annotate(postgres.NewDBX, fx.As(new(gophermart.Querier))),
			fx.Annotate(accrual.NewClient, fx.As(new(gophermart.AccrualClient))),
			fx.Annotate(auth.NewService, fx.As(new(gophermart.AuthService))),
			fx.Annotate(withdraw.NewService, fx.As(new(gophermart.WithdrawService))),
			fx.Annotate(orders.NewService, fx.As(new(gophermart.OrderService))),
			NewHTTPService,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, conn *pgx.Conn, logger *zap.SugaredLogger) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						logger.Info("running database DDL migrations...")
						return runMigrations(ctx, conn)
					},
				})
			},
			func(lc fx.Lifecycle, ordersService gophermart.OrderService) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						ordersService.StartAccrualFetching(ctx)
						return nil
					},
				})
			},
			func(*http.Server) {},
		),
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
