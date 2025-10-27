package ship_assembled_producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/Daniil-Sakharov/RocketFactory/assembly/internal/model"
	def "github.com/Daniil-Sakharov/RocketFactory/assembly/internal/service"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/kafka"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
	eventsv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/events/v1"
)

var _ def.ShipAssembledProducerService = (*service)(nil)

type service struct {
	shipAssembledProducer kafka.Producer
}

func NewService(shipAssembledProducer kafka.Producer) *service {
	return &service{
		shipAssembledProducer: shipAssembledProducer,
	}
}

func (s *service) PublishShipAssembled(ctx context.Context, event *model.ShipAssembledEvent) error {
	msg := &eventsv1.ShipAssembled{
		EventUuid:    event.EventUUID,
		OrderUuid:    event.OrderUUID,
		UserUuid:     event.UserUUID,
		BuildTimeSec: int64(event.BuildTime.Seconds()),
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "Failed to marshal ShipAssembled event", zap.Error(err))
		return err
	}

	err = s.shipAssembledProducer.Send(ctx, []byte(event.OrderUUID), payload)
	if err != nil {
		logger.Error(ctx, "Failed to publish ShipAssembled event", zap.Error(err))
		return err
	}

	logger.Info(ctx, "📤 ShipAssembled event published",
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.Int64("build_time_sec", int64(event.BuildTime.Seconds())),
	)

	return nil
}
