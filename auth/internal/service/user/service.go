package user

import (
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/repository"
	def "github.com/Daniil-Sakharov/RocketFactory/auth/internal/service"
)

var _ def.UserService = (*service)(nil)

type service struct {
	userRepository    repository.UsersRepository
	sessionRepository repository.SessionCacheRepository
}

func NewService(
	userRepository repository.UsersRepository,
	sessionRepository repository.SessionCacheRepository,
) *service {
	return &service{
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
	}
}
