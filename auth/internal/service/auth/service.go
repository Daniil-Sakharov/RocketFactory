package auth

import (
	"time"

	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/repository"
	def "github.com/Daniil-Sakharov/RocketFactory/auth/internal/service"
)

var _ def.AuthService = (*service)(nil)

type service struct {
	userRepository    repository.UsersRepository
	sessionRepository repository.SessionCacheRepository
	sessionTTL        time.Duration
}

func NewService(
	userRepository repository.UsersRepository,
	sessionRepository repository.SessionCacheRepository,
	sessionTTL time.Duration,
) *service {
	return &service{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
		sessionTTL:        sessionTTL,
	}
}
