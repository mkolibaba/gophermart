package http

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/mkolibaba/gophermart"
	"github.com/mkolibaba/gophermart/internal/config"
	"github.com/mkolibaba/gophermart/internal/http/handler/auth"
	"github.com/mkolibaba/gophermart/internal/http/handler/balance"
	"github.com/mkolibaba/gophermart/internal/http/handler/orders"
	"github.com/mkolibaba/gophermart/internal/http/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Server struct {
	Address         string
	Logger          *zap.SugaredLogger
	Querier         gophermart.Querier
	AuthService     gophermart.AuthService
	OrderService    gophermart.OrderService
	WithdrawService gophermart.WithdrawService

	router *echo.Echo
}

func NewServer(
	lc fx.Lifecycle,
	cfg *config.Config,
	logger *zap.SugaredLogger,
	dbx gophermart.Querier,
	authService gophermart.AuthService,
	ordersService gophermart.OrderService,
	withdrawService gophermart.WithdrawService,
) *Server {
	h := &Server{
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
		OnStop: h.Stop,
	})

	return h
}

func (s *Server) Start(ctx context.Context) {
	s.router = echo.New()

	// handlers
	authHandler := auth.NewHandler(s.Querier, s.AuthService)
	ordersHandler := orders.NewHandler(s.OrderService)
	balanceHandler := balance.NewHandler(s.Querier, s.WithdrawService)

	apiUserRouter := s.router.Group("/api/user")
	apiUserRouter.POST("/register", authHandler.Register)
	apiUserRouter.POST("/login", authHandler.Login)

	securedRouter := apiUserRouter.Group("")
	securedRouter.Use(middleware.Auth(s.AuthService))

	securedRouter.POST("/orders", ordersHandler.Create)
	securedRouter.GET("/orders", ordersHandler.GetAll)
	securedRouter.GET("/balance", balanceHandler.Get)
	securedRouter.POST("/balance/withdraw", balanceHandler.WithdrawBalance)
	securedRouter.GET("/withdrawals", balanceHandler.GetAllWithdrawals)

	s.router.Start(s.Address)
}

func (s *Server) Stop(ctx context.Context) error {
	return s.router.Shutdown(ctx)
}
