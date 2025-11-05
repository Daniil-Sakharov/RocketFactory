package v1

import (
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/service"
	pb "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/auth/v1"
)

var _ pb.AuthServiceServer = (*api)(nil)

type api struct {
	pb.UnimplementedAuthServiceServer
	authService service.AuthService
}

func NewAPI(authService service.AuthService) *api {
	return &api{
		authService: authService,
	}
}
