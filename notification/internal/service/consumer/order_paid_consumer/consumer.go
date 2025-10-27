package order_paid_consumer

import (
	"context"
	kafkaConverter "github.com/Daniil-Sakharov/RocketFactory/notification/internal/converter/kafka"
	serv "github.com/Daniil-Sakharov/RocketFactory/notification/internal/service"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/kafka"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
	"go.uber.org/zap"
)

var _ serv.OrderPaidConsumerService = (*service)(nil)

type service struct {
	orderPaidConsumer kafka.Consumer
	orderPaidDecoder  kafkaConverter.OrderDecoder
	telegramService      serv.TelegramService
}

func NewService(
	orderPaidConsumer kafka.Consumer,
	orderPaidDecoder kafkaConverter.OrderDecoder,
	telegramService serv.TelegramService,
) *service {
	return &service{
		orderPaidConsumer: orderPaidConsumer,
		orderPaidDecoder:  orderPaidDecoder,
		telegramService:   telegramService,
	}
}

func (s *service) RunOrderConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 Starting ShipAssembly consumer service")

	err := s.orderPaidConsumer.Consume(ctx, s.handleOrderPaid)
	if err != nil {
		logger.Error(ctx, "❌ Failed to consume from ship.assembled topic", zap.Error(err))
		return err
	}

	return nil
}
