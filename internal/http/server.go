package http

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/mkolibaba/gophermart"
	"github.com/mkolibaba/gophermart/internal/http/handler/auth"
	"github.com/mkolibaba/gophermart/internal/http/handler/balance"
	"github.com/mkolibaba/gophermart/internal/http/handler/orders"
	"github.com/mkolibaba/gophermart/internal/http/middleware"
	"go.uber.org/zap"
)

type Server struct {
	Address         string
	Logger          *zap.SugaredLogger
	Querier         gophermart.Querier
	OrderService    gophermart.OrderService
	WithdrawService gophermart.WithdrawService
}

func (s *Server) Start(ctx context.Context) {
	router := echo.New()

	// handlers
	authHandler := auth.NewHandler(s.Querier)
	ordersHandler := orders.NewHandler(s.OrderService)
	balanceHandler := balance.NewHandler(s.Querier, s.WithdrawService)

	apiUserRouter := router.Group("/api/user")
	apiUserRouter.POST("/register", authHandler.Register)
	apiUserRouter.POST("/login", authHandler.Login)

	securedRouter := apiUserRouter.Group("")
	securedRouter.Use(middleware.Auth())

	securedRouter.POST("/orders", ordersHandler.Create)
	securedRouter.GET("/orders", ordersHandler.GetAll)
	securedRouter.GET("/balance", balanceHandler.Get)
	securedRouter.POST("/balance/withdraw", balanceHandler.WithdrawBalance)
	securedRouter.GET("/withdrawals", balanceHandler.GetAllWithdrawals)

	// TODO(improvement): graceful shutdown
	router.Start(s.Address)
}
