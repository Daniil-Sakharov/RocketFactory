package session

import (
	def "github.com/Daniil-Sakharov/RocketFactory/auth/internal/repository"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/cache"
)

var _ def.SessionCacheRepository = (*repository)(nil)

type repository struct {
	cache cache.RedisClient
}

func NewRepository(cache cache.RedisClient) *repository {
	return &repository{cache: cache}
}
