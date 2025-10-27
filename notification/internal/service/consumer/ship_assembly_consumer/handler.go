package ship_assembly_consumer

import (
	"context"
	"errors"
	converter "github.com/Daniil-Sakharov/RocketFactory/notification/internal/converter/telegram"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/kafka/consumer"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) handleShipAssembly(ctx context.Context, msg consumer.Message) error {
	event, err := s.shipAssemblyDecoder.AssemblyDecode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode ShipAssembly event", zap.Error(err))
		return err
	}
	if event.EventUUID == "" {
		logger.Error(ctx, "Invalid event: empty order_uuid")
		return errors.New("invalid event")
	}

	logger.Info(ctx, "📨 Received ShipAssembled event",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("build_time_sec", event.BuildTimeSec.String()),
	)

	err = s.telegramService.SendShipAssembledNotification(ctx, converter.ShipAssembledEventToTemplateData(&event))
	if err != nil {
		logger.Error(ctx, "Failed to send ShipAssembly event to telegram", zap.Error(err))
		return err
	}

	logger.Info(ctx, "✅ ShipAssembly event processed successfully",
		zap.String("order_uuid", event.OrderUUID),
	)

	return nil
}
