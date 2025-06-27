package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/mkolibaba/gophermart/internal/auth"
	"github.com/mkolibaba/gophermart/internal/config"
	"github.com/mkolibaba/gophermart/internal/http"
	"github.com/mkolibaba/gophermart/internal/http/client/accrual"
	"github.com/mkolibaba/gophermart/internal/orders"
	"github.com/mkolibaba/gophermart/internal/withdraw"
	"github.com/mkolibaba/gophermart/postgres"
	"github.com/mkolibaba/gophermart/postgres/migration"
	"go.uber.org/zap"
	stdlog "log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		stdlog.Fatal(err)
	}

	unsugaredLogger, err := zap.NewDevelopment()
	if err != nil {
		stdlog.Fatal(err)
	}
	logger := unsugaredLogger.Sugar()

	logger.Infof("provided configuration: %+v", cfg)

	conn, err := pgx.Connect(ctx, cfg.DatabaseURI)
	if err != nil {
		logger.Fatal(err)
	}

	dbx := postgres.NewDBX(conn)

	// TODO(improvement): использовать систему миграции
	logger.Info("running database DDL migrations...")
	if _, err := conn.Exec(ctx, migration.DDL); err != nil {
		logger.Fatalf("failed to run database ddl migrations: %s", err)
	}

	accrualClient := accrual.NewClient(cfg.AccrualSystemAddress, logger)
	authService := auth.NewService()
	withdrawService := withdraw.NewService(dbx, dbx)
	ordersService := orders.NewService(dbx, accrualClient, logger)
	ordersService.StartAccrualFetching(ctx)

	server := &http.Server{
		Address:         cfg.RunAddress,
		Logger:          logger,
		Querier:         dbx,
		AuthService:     authService,
		OrderService:    ordersService,
		WithdrawService: withdrawService,
	}
	server.Start(ctx)
}
