package session

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	repo "github.com/Daniil-Sakharov/RocketFactory/auth/internal/repository"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
)

func (r *repository) Delete(ctx context.Context, sessionUUID, userUUID string) error {
	sessionKey := repo.SessionKey(sessionUUID)

	if err := r.cache.Del(ctx, sessionKey); err != nil {
		return fmt.Errorf("failed to delete session from redis: %w", err)
	}

	if err := r.removeSessionFromUserSet(ctx, userUUID, sessionUUID); err != nil {
		logger.Warn(ctx, "failed to remove session from user set", zap.Error(err))
	}

	return nil
}

func (r *repository) removeSessionFromUserSet(ctx context.Context, userUUID, sessionUUID string) error {
	userSessionKey := repo.UserSessionKey(userUUID)

	if err := r.cache.SRem(ctx, userSessionKey, sessionUUID); err != nil {
		return fmt.Errorf("failed to remove session from user set: %w", err)
	}

	return nil
}
