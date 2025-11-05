package v1

import (
	"context"

	"go.uber.org/zap"

	converter2 "github.com/Daniil-Sakharov/RocketFactory/auth/internal/converter"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
	authv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/auth/v1"
)

func (a *api) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	dtoReq := converter2.LoginRequestToDTO(req)

	dtoResp, err := a.authService.Login(ctx, dtoReq)
	if err != nil {
		logger.Warn(ctx, "login request failed",
			zap.String("login", req.Login),
			zap.Error(err),
		)

		return nil, converter2.MapErrorToGRPC(err)
	}

	return converter2.LoginResponseFromDTO(dtoResp), nil
}
