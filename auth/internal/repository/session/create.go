package session

import (
	"context"
	"fmt"
	"time"

	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model/entity"
	repo "github.com/Daniil-Sakharov/RocketFactory/auth/internal/repository"
	repoConverter "github.com/Daniil-Sakharov/RocketFactory/auth/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, session *entity.Session, ttl time.Duration) error {
	sessionRepo := repoConverter.EntityToRepositorySession(session)

	sessionKey := repo.SessionKey(sessionRepo.SessionUUID)

	redisView := repoConverter.SessionToRedisView(sessionRepo)

	if err := r.cache.HashSet(ctx, sessionKey, redisView); err != nil {
		return fmt.Errorf("failed to save session to redis hash: %w", err)
	}

	if err := r.cache.Expire(ctx, sessionKey, ttl); err != nil {
		//nolint:gosec // Best-effort cleanup on error path
		_ = r.cache.Del(ctx, sessionKey)
		return fmt.Errorf("failed to set ttl for session: %w", err)
	}

	if err := r.addSessionToUserSet(ctx, sessionRepo.UserUUID, sessionRepo.SessionUUID); err != nil {
		//nolint:gosec // Best-effort cleanup on error path
		_ = r.cache.Del(ctx, sessionKey)
		return fmt.Errorf("failed to add session to user set: %w", err)
	}

	return nil
}

func (r *repository) addSessionToUserSet(ctx context.Context, userUUID, sessionUUID string) error {
	userSessionKey := repo.UserSessionKey(userUUID)

	if err := r.cache.SAdd(ctx, userSessionKey, sessionUUID); err != nil {
		return fmt.Errorf("failed to add session to user set: %w", err)
	}

	return nil
}
