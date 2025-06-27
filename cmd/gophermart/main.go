package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/gommon/log"
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
	ctx := context.Background()

	cfg, err := config.New()
	if err != nil {
		fail(err)
	}

	unsugaredLogger, err := zap.NewDevelopment()
	if err != nil {
		fail(err)
	}
	logger := unsugaredLogger.Sugar()

	conn, err := pgx.Connect(ctx, cfg.DatabaseURI)
	if err != nil {
		fail(err)
	}

	dbx := postgres.NewDBX(conn)

	// TODO(improvement): использовать систему миграции
	logger.Info("running database DDL migrations...")
	if _, err := conn.Exec(ctx, migration.DDL); err != nil {
		log.Fatalf("failed to run database ddl migrations: %s", err)
	}

	accrualClient := accrual.NewClient(cfg.AccrualSystemAddress, logger)
	withdrawService := withdraw.NewService(dbx, dbx)
	ordersService := orders.NewService(dbx, accrualClient, logger)
	ordersService.StartAccrualFetching(ctx)

	server := &http.Server{
		Address:         cfg.RunAddress,
		Logger:          logger,
		Querier:         dbx,
		OrderService:    ordersService,
		WithdrawService: withdrawService,
	}
	server.Start(ctx)
}

func fail(err error) {
	stdlog.Fatal(err)
}
