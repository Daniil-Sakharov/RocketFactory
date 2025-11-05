package v1

import (
	"context"

	"go.uber.org/zap"

	converter2 "github.com/Daniil-Sakharov/RocketFactory/auth/internal/converter"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
	authv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/auth/v1"
)

func (a *api) Whoami(ctx context.Context, req *authv1.WhoamiRequest) (*authv1.WhoamiResponse, error) {
	dtoReq := converter2.WhoamiRequestToDTO(req)

	dtoResp, err := a.authService.Whoami(ctx, dtoReq)
	if err != nil {
		logger.Warn(ctx, "whoami request failed",
			zap.String("session_uuid", req.SessionUuid),
			zap.Error(err),
		)

		return nil, converter2.MapErrorToGRPC(err)
	}

	return converter2.WhoamiResponseFromDTO(dtoResp), nil
}
