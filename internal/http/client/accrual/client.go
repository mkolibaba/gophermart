package accrual

import (
	"context"
	"github.com/mkolibaba/gophermart"
	"go.uber.org/zap"
	"resty.dev/v3"
)

type Client struct {
	client *resty.Client
}

func NewClient(accrualAddress string, logger *zap.SugaredLogger) *Client {
	return &Client{
		client: resty.New().
			SetBaseURL(accrualAddress).
			SetScheme("http").
			SetRetryCount(2).
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

	// TODO: обработать возможные статусы
	return response.Result().(*gophermart.AccrualOrder), nil
}
