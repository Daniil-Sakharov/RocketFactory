package assembly_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/Daniil-Sakharov/RocketFactory/order/internal/converter/kafka"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/repository"
	serv "github.com/Daniil-Sakharov/RocketFactory/order/internal/service"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/kafka"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
)

type service struct {
	orderConsumer   kafka.Consumer
	orderDecoder    kafkaConverter.AssemblyDecoder
	orderService    serv.OrderService
	orderRepository repository.OrderRepository
}

func NewService(
	orderConsumer kafka.Consumer,
	orderDecoder kafkaConverter.AssemblyDecoder,
	orderService serv.OrderService,
	orderRepository repository.OrderRepository,
) *service {
	return &service{
		orderConsumer:   orderConsumer,
		orderDecoder:    orderDecoder,
		orderService:    orderService,
		orderRepository: orderRepository,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting ShipAssembled service")

	err := s.orderConsumer.Consume(ctx, s.orderHandler)
	if err != nil {
		logger.Error(ctx, "Failed to consume from ship.assembled topic", zap.Error(err))
		return err
	}

	return nil
}
