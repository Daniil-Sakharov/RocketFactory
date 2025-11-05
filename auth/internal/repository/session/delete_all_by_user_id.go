package session

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	repo "github.com/Daniil-Sakharov/RocketFactory/auth/internal/repository"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
)

func (r *repository) DeleteAllByUserUUID(ctx context.Context, userUUID string) error {
	userSessionKey := repo.UserSessionKey(userUUID)

	sessionUUIDs, err := r.cache.SMembers(ctx, userSessionKey)
	if err != nil {
		return fmt.Errorf("failed to get user session from redis: %w", err)
	}

	for _, sessionUUID := range sessionUUIDs {
		sessionKey := repo.SessionKey(sessionUUID)
		if err = r.cache.Del(ctx, sessionKey); err != nil {
			logger.Error(ctx, "failed to delete session",
				zap.String("session_uuid", sessionUUID),
				zap.Error(err),
			)
		}
	}

	if err = r.cache.Del(ctx, userSessionKey); err != nil {
		return fmt.Errorf("failed to delete user sessions set from redis: %w", err)
	}

	return nil
}
