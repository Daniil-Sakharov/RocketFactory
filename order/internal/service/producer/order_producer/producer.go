package order_producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/domain"
	def "github.com/Daniil-Sakharov/RocketFactory/order/internal/service"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/kafka"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
	eventsv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/events/v1"
)

var _ def.OrderProducerService = (*service)(nil)

type service struct {
	orderProducer kafka.Producer
}

func NewService(orderProducer kafka.Producer) *service {
	return &service{
		orderProducer: orderProducer,
	}
}

func (s *service) PublishOrderPaid(ctx context.Context, event *domain.OrderProduceEvent) error {
	msg := &eventsv1.OrderPaid{
		EventUuid:       event.EventUUID,
		OrderUuid:       event.OrderUUID,
		UserUuid:        event.UserUUID,
		PaymentMethod:   event.PaymentMethod,
		TransactionUuid: event.TransactionUUID,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "Failed to marshal ShipAssembled event", zap.Error(err))
		return err
	}

	err = s.orderProducer.Send(ctx, []byte(event.OrderUUID), payload)
	if err != nil {
		logger.Error(ctx, "Failed to publish OrderPaid event", zap.Error(err))
		return err
	}

	logger.Info(ctx, "📤 OrderPaid event published",
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.String("payment_method", event.PaymentMethod),
		zap.String("transaction_uuid", event.TransactionUUID),
	)

	return nil
}
