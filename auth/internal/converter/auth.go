package converter

import (
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/service/dto"
	pbAuth "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/auth/v1"
	pbCommon "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/common/v1"
)

func LoginRequestToDTO(req *pbAuth.LoginRequest) *dto.LoginRequest {
	return &dto.LoginRequest{
		Login:    req.Login,
		Password: req.Password,
	}
}

func LoginResponseFromDTO(resp *dto.LoginResponse) *pbAuth.LoginResponse {
	return &pbAuth.LoginResponse{
		SessionUuid: resp.SessionUUID,
	}
}

func WhoamiRequestToDTO(req *pbAuth.WhoamiRequest) *dto.WhoamiRequest {
	return &dto.WhoamiRequest{
		SessionUUID: req.SessionUuid,
	}
}

func WhoamiResponseFromDTO(resp *dto.WhoamiResponse) *pbAuth.WhoamiResponse {
	return &pbAuth.WhoamiResponse{
		User: &pbCommon.User{
			UserUuid: resp.UserUUID,
			Login:    resp.Login,
			Email:    resp.Email,
		},
	}
}
