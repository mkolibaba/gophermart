package accrual

import (
	"context"
	"fmt"
	"github.com/mkolibaba/gophermart"
	"go.uber.org/zap"
	"net/http"
	"resty.dev/v3"
	"time"
)

type Client struct {
	client *resty.Client
}

func NewClient(accrualAddress string, logger *zap.SugaredLogger) *Client {
	return &Client{
		client: resty.New().
			SetBaseURL(accrualAddress).
			SetScheme("http").
			// retry настройки учитывают 429 (и Retry-After хедер) и 500+ статусы
			SetRetryCount(2).
			SetRetryWaitTime(2 * time.Second).
			SetRetryMaxWaitTime(5 * time.Second).
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

	if response.StatusCode() == http.StatusNoContent {
		return nil, fmt.Errorf("order has not been registered")
	}

	return response.Result().(*gophermart.AccrualOrder), nil
}
