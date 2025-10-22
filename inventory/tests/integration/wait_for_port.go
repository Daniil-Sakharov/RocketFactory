//go:build integration

package integration

import (
	"context"
	"net"
	"time"

	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
	"go.uber.org/zap"
)

// waitForPort waits for a TCP port to become available
func waitForPort(ctx context.Context, address string, maxAttempts int, delay time.Duration) error {
	logger.Info(ctx, "⏳ Waiting for port to be available", zap.String("address", address))

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		conn, err := net.DialTimeout("tcp", address, 2*time.Second)
		if err == nil {
			conn.Close()
			logger.Info(ctx, "✅ Port is now available",
				zap.String("address", address),
				zap.Int("attempt", attempt))
			return nil
		}

		logger.Debug(ctx, "Port not yet available",
			zap.String("address", address),
			zap.Int("attempt", attempt),
			zap.Int("maxAttempts", maxAttempts),
			zap.Error(err))

		if attempt < maxAttempts {
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}

	return logger.Error(ctx, "Port did not become available",
		zap.String("address", address),
		zap.Int("maxAttempts", maxAttempts))
}
