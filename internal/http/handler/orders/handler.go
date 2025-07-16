package orders

import (
	"github.com/mkolibaba/gophermart"
)

type Handler struct {
	orderService gophermart.OrderService
}

func NewHandler(orderService gophermart.OrderService) *Handler {
	return &Handler{
		orderService: orderService,
	}
}
