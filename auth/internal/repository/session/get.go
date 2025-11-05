package session

import (
	"context"
	"errors"
	"fmt"

	redigo "github.com/gomodule/redigo/redis"

	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model/entity"
	repo "github.com/Daniil-Sakharov/RocketFactory/auth/internal/repository"
	repoConverter "github.com/Daniil-Sakharov/RocketFactory/auth/internal/repository/converter"
	repoModel "github.com/Daniil-Sakharov/RocketFactory/auth/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, sessionUUID string) (*entity.Session, error) {
	sessionKey := repo.SessionKey(sessionUUID)

	values, err := r.cache.HGetAll(ctx, sessionKey)
	if err != nil {
		if errors.Is(err, redigo.ErrNil) {
			return nil, model.ErrSessionNotFound
		}
		return nil, fmt.Errorf("failed to get session from redis: %w", err)
	}

	if len(values) == 0 {
		return nil, model.ErrSessionNotFound
	}

	var sessionRedisView repoModel.SessionRedisView
	if err = redigo.ScanStruct(values, &sessionRedisView); err != nil {
		return nil, fmt.Errorf("failed to scan sesson from redis hash: %w", err)
	}

	sessionRepo := repoConverter.SessionFromRedisView(&sessionRedisView)

	session := repoConverter.RepositoryToEntitySession(sessionRepo)

	return session, nil
}
