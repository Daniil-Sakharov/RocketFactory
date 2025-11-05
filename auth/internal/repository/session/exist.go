package session

import (
	"context"
	"fmt"

	repo "github.com/Daniil-Sakharov/RocketFactory/auth/internal/repository"
)

func (r *repository) Exist(ctx context.Context, sessionUUID string) (bool, error) {
	sessionKey := repo.SessionKey(sessionUUID)

	exist, err := r.cache.Exists(ctx, sessionKey)
	if err != nil {
		return false, fmt.Errorf("failed to check session existence in redis: %w", err)
	}
	return exist, nil
}
