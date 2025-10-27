package assembly

import (
	"context"
	"crypto/rand"
	"math/big"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/RocketFactory/assembly/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
)

func (s *service) ProcessOrderPaid(ctx context.Context, event *model.OrderPaidEvent) error {
	logger.Info(ctx, "🚀 Starting ship assembly",
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
	)

	randomSeconds, err := rand.Int(rand.Reader, big.NewInt(10))
	if err != nil {
		logger.Error(ctx, "Failed to generate random build time", zap.Error(err))
		return err
	}
	buildTime := time.Duration(randomSeconds.Int64()+1) * time.Second

	logger.Info(ctx, "⏳ Assembling ship...",
		zap.String("order_uuid", event.OrderUUID),
		zap.Int("build_time_sec", int(buildTime.Seconds())),
	)

	timer := time.NewTimer(buildTime)
	select {
	case <-ctx.Done():
		timer.Stop()
		return ctx.Err()
	case <-timer.C:
	}

	logger.Info(ctx, "✅ Ship assembled successfully",
		zap.String("order_uuid", event.OrderUUID),
		zap.Int("build_time_sec", int(buildTime.Seconds())),
	)

	shipAssembledEvent := &model.ShipAssembledEvent{
		EventUUID: uuid.New().String(),
		OrderUUID: event.OrderUUID,
		UserUUID:  event.UserUUID,
		BuildTime: buildTime,
	}

	if err = s.shipAssembledProducer.PublishShipAssembled(ctx, shipAssembledEvent); err != nil {
		logger.Error(ctx, "Failed to publish ShipAssembled event", zap.Error(err))
		return err
	}

	return nil
}
