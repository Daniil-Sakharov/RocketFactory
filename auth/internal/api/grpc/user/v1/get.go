package v1

import (
	"context"

	"go.uber.org/zap"

	converter2 "github.com/Daniil-Sakharov/RocketFactory/auth/internal/converter"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
	userv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/user/v1"
)

func (a *api) Get(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	dtoReq := converter2.GetUserRequestToDTO(req)

	dtoResp, err := a.userService.Get(ctx, dtoReq)
	if err != nil {
		logger.Warn(ctx, "get request failed",
			zap.String("user_uuid", req.UserUuid),
		)
		return nil, converter2.MapErrorToGRPC(err)
	}

	return converter2.GetUserResponseFromDTO(dtoResp), nil
}
