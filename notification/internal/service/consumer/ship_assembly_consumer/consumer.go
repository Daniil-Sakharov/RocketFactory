package ship_assembly_consumer

import (
	"context"
	kafkaConverter "github.com/Daniil-Sakharov/RocketFactory/notification/internal/converter/kafka"
	serv "github.com/Daniil-Sakharov/RocketFactory/notification/internal/service"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/kafka"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
	"go.uber.org/zap"
)

var _ serv.ShipAssemblyConsumerService = (*service)(nil)

type service struct {
	shipAssemblyConsumer kafka.Consumer
	shipAssemblyDecoder  kafkaConverter.AssemblyDecoder
	telegramService      serv.TelegramService
}

func NewService(
	shipAssemblyConsumer kafka.Consumer,
	shipAssemblyDecoder kafkaConverter.AssemblyDecoder,
	telegramService serv.TelegramService,
) *service {
	return &service{
		shipAssemblyConsumer: shipAssemblyConsumer,
		shipAssemblyDecoder:  shipAssemblyDecoder,
		telegramService:      telegramService,
	}
}

func (s *service) RunAssemblyConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 Starting ShipAssembly consumer service")

	err := s.shipAssemblyConsumer.Consume(ctx, s.handleShipAssembly)
	if err != nil {
		logger.Error(ctx, "❌ Failed to consume from ship.assembled topic", zap.Error(err))
		return err
	}

	return nil
}
