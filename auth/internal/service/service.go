package service

import (
	"context"

	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/service/dto"
)

type UserService interface {
	Register(ctx context.Context, req *dto.RegisterUserRequest) (*dto.RegisterUserResponse, error)
	Get(ctx context.Context, req *dto.GetUserRequest) (*dto.GetUserResponse, error)
	GetByLogin(ctx context.Context, login string) (*dto.GetUserResponse, error)
	Update(ctx context.Context, req *dto.UpdateUserRequest) error
}

type AuthService interface {
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	Whoami(ctx context.Context, req *dto.WhoamiRequest) (*dto.WhoamiResponse, error)
}
