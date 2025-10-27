package assembly

import (
	"context"
	"math/rand"
	"time"

	"github.com/Daniil-Sakharov/RocketFactory/assembly/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *service) ProcessOrderPaid(ctx context.Context, event *model.OrderPaidEvent) error {
	logger.Info(ctx, "🚀 Starting ship assembly",
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
	)

	buildTime := time.Duration(rand.Intn(10)+1) * time.Second

	logger.Info(ctx, "⏳ Assembling ship...",
		zap.String("order_uuid", event.OrderUUID),
		zap.Int("build_time_sec", int(buildTime.Seconds())),
	)

	time.Sleep(buildTime)

	logger.Info(ctx, "✅ Ship assembled successfully",
		zap.String("order_uuid", event.OrderUUID),
		zap.Int("build_time_sec", int(buildTime.Seconds())),
	)

	shipAssembledEvent := &model.ShipAssembledEvent{
		EventUUID:    uuid.New().String(),
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: buildTime,
	}

	err := s.shipAssembledProducer.PublishShipAssembled(ctx, shipAssembledEvent)
	if err != nil {
		logger.Error(ctx, "Failed to publish ShipAssembled event", zap.Error(err))
		return err
	}

	return nil
}
