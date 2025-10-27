package assembly_consumer

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/vo"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/kafka/consumer"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
)

func (s *service) orderHandler(ctx context.Context, msg consumer.Message) error {
	event, err := s.orderDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode ShipAssembled event")
		return err
	}

	if event.OrderUUID == "" {
		logger.Error(ctx, "Invalid event: empty order_uuid")
		return errors.New("invalid event")
	}

	logger.Info(ctx, "📨 Received OrderPaid event",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.Int("build_time_sec", int(event.BuildTime.Seconds())),
	)

	order, err := s.orderRepository.Get(ctx, event.OrderUUID)
	if err != nil {
		logger.Error(ctx, "Failed to get order", zap.Error(err))
		return err
	}

	order.Status = vo.OrderStatusASSEMBLED

	err = s.orderRepository.Update(ctx, order)
	if err != nil {
		logger.Error(ctx, "Failed to update order status to ASSEMBLED", zap.Error(err))
		return err
	}

	logger.Info(ctx, "✅ Order status updated to ASSEMBLED",
		zap.String("order_uuid", event.OrderUUID),
	)

	return nil
}
