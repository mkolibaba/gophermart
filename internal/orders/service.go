package orders

import (
	"context"
	"fmt"
	"github.com/mkolibaba/gophermart"
	"github.com/mkolibaba/gophermart/postgres/gen"
	"go.uber.org/zap"
	"time"
)

const (
	workers         = 5
	fetchTick       = 15 * time.Second
	ordersBatchSize = workers
)

var accrualStatusMapping = map[string]postgres.AccrualStatus{
	"REGISTERED": postgres.AccrualStatusPROCESSING,
	"INVALID":    postgres.AccrualStatusINVALID,
	"PROCESSING": postgres.AccrualStatusPROCESSING,
	"PROCESSED":  postgres.AccrualStatusPROCESSED,
}

type Service struct {
	gophermart.Querier
	accrualClient gophermart.AccrualClient
	logger        *zap.SugaredLogger
}

func NewService(querier gophermart.Querier, accrualClient gophermart.AccrualClient, logger *zap.SugaredLogger) *Service {
	return &Service{
		Querier:       querier,
		accrualClient: accrualClient,
		logger:        logger,
	}
}

func (s *Service) StartAccrualFetching(ctx context.Context) {
	ticker := time.NewTicker(fetchTick)

	worker := func(in <-chan *postgres.Order) {
		for order := range in {
			response, err := s.accrualClient.GetOrder(ctx, order.ID)
			if err != nil {
				s.logger.Errorf("failed to get accrual response: %s", err)
				return
			}

			if status, ok := accrualStatusMapping[response.Status]; ok {
				err := s.DoInTx(ctx, func(qtx postgres.Querier) error {
					if err := qtx.OrderUpdateAccrualStatus(ctx, status, order.ID); err != nil {
						return fmt.Errorf("accrual status update: %w", err)
					}

					if response.Status == "PROCESSED" {
						if err := qtx.UserAddAccrualBalance(ctx, response.Accrual, order.UserLogin); err != nil {
							return fmt.Errorf("accrual balance update: %w", err)
						}
					}

					return nil
				})
				if err != nil {
					s.logger.Error(err)
				}
			}
		}
	}

	// TODO(improvement): никто не гарантирует, что за тик обработаются все заказы,
	//  т.е. в канале может оказаться дважды один и тот же заказ
	ch := make(chan *postgres.Order, ordersBatchSize)

	for range workers {
		go worker(ch)
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				result, err := s.OrderGetWithNonFinalAccrualStatus(ctx, ordersBatchSize)
				if err != nil {
					s.logger.Errorf("failed to get orders to process: %s", err)
				}
				s.logger.Debugf("there are %d orders to process", len(result))
				for _, r := range result {
					ch <- r
				}
			case <-ctx.Done():
				ticker.Stop()
				close(ch)
				return
			}
		}
	}()
}
