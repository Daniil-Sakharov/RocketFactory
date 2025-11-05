package repository

import (
	"context"
	"time"

	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model/entity"
)

type UsersRepository interface {
	Create(ctx context.Context, users *entity.User) error
	Get(ctx context.Context, uuid string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	GetByLogin(ctx context.Context, login string) (*entity.User, error)
}

type SessionCacheRepository interface {
	Create(ctx context.Context, session *entity.Session, ttl time.Duration) error
	Get(ctx context.Context, sessionUUID string) (*entity.Session, error)
	Delete(ctx context.Context, sessionUUID, userUUID string) error
	DeleteAllByUserUUID(ctx context.Context, userUUID string) error
	Exist(ctx context.Context, sessionUUID string) (bool, error)
}
