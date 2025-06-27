package orders

import (
	"context"
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
	defer ticker.Stop()

	worker := func(in <-chan *postgres.Order) {
		for order := range in {
			response, err := s.accrualClient.GetOrder(ctx, order.ID)
			if err != nil {
				s.logger.Errorf("failed to get accrual response: %s", err)
				return
			}

			if status, ok := accrualStatusMapping[response.Status]; ok {
				if err := s.OrderUpdateAccrualStatus(ctx, status, order.ID); err != nil {
					s.logger.Errorf("failed to update order accrual status: %s", err)
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
				for _, r := range result {
					ch <- r
				}
			case <-ctx.Done():
				close(ch)
				return
			}
		}
	}()
}
