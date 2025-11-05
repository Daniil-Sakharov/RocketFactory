package v1

import (
	"context"

	"go.uber.org/zap"

	converter2 "github.com/Daniil-Sakharov/RocketFactory/auth/internal/converter"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
	userv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/user/v1"
)

func (a *api) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	reqDto := converter2.RegisterRequestToDTO(req)

	respDto, err := a.userService.Register(ctx, reqDto)
	if err != nil {
		logger.Warn(ctx, "register request failed",
			zap.String("login", req.Login),
			zap.String("email", req.Email),
		)

		return nil, converter2.MapErrorToGRPC(err)
	}

	return converter2.RegisterResponseFromDTO(respDto), nil
}
