package accrual

import (
	"context"
	"fmt"
	"github.com/mkolibaba/gophermart"
	"go.uber.org/zap"
	"net/http"
	"resty.dev/v3"
	"strconv"
)

type Client struct {
	client *resty.Client
}

func NewClient(accrualAddress string, logger *zap.SugaredLogger) *Client {
	return &Client{
		client: resty.New().
			SetBaseURL(accrualAddress).
			SetScheme("http").
			SetLogger(logger),
	}
}

func (c *Client) GetOrder(ctx context.Context, number string) (*gophermart.AccrualOrder, error) {
	response, err := c.client.R().
		SetContext(ctx).
		SetResult(&gophermart.AccrualOrder{}).
		Get("/api/orders/" + number)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode() {
	case http.StatusNoContent:
		return nil, fmt.Errorf("order has not been registered")
	case http.StatusTooManyRequests:
		interval, err := strconv.Atoi(response.Header().Get("Retry-After"))
		if err != nil {
			return nil, err
		}
		return nil, gophermart.RetryAfterError{Interval: int64(interval)}
	default:
		return response.Result().(*gophermart.AccrualOrder), nil
	}
}
