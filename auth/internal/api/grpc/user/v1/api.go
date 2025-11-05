package v1

import (
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/service"
	pb "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/user/v1"
)

var _ pb.UserServiceServer = (*api)(nil)

type api struct {
	pb.UnimplementedUserServiceServer
	userService service.UserService
}

func NewAPI(userService service.UserService) *api {
	return &api{
		userService: userService,
	}
}
