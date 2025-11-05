package converter

import (
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/service/dto"
	pbCommon "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/common/v1"
	pbUser "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/user/v1"
)

func RegisterRequestToDTO(req *pbUser.RegisterRequest) *dto.RegisterUserRequest {
	notificationMethods := make([]dto.NotificationMethodDTO, 0, len(req.NotificationMethods))
	for _, nm := range req.NotificationMethods {
		notificationMethods = append(notificationMethods, dto.NotificationMethodDTO{
			ProviderName: nm.ProviderName,
			Target:       nm.Target,
		})
	}

	return &dto.RegisterUserRequest{
		Login:               req.Login,
		Password:            req.Password,
		Email:               req.Email,
		NotificationMethods: notificationMethods,
	}
}

func RegisterResponseFromDTO(resp *dto.RegisterUserResponse) *pbUser.RegisterResponse {
	return &pbUser.RegisterResponse{
		UserUuid: resp.UserUUID,
	}
}

func GetUserRequestToDTO(req *pbUser.GetUserRequest) *dto.GetUserRequest {
	return &dto.GetUserRequest{
		UserUUID: req.UserUuid,
	}
}

func GetUserResponseFromDTO(resp *dto.GetUserResponse) *pbUser.GetUserResponse {
	notificationMethods := make([]*pbCommon.NotificationMethod, 0, len(resp.NotificationMethods))
	for _, nm := range resp.NotificationMethods {
		notificationMethods = append(notificationMethods, &pbCommon.NotificationMethod{
			ProviderName: nm.ProviderName,
			Target:       nm.Target,
		})
	}

	return &pbUser.GetUserResponse{
		User: &pbCommon.User{
			UserUuid: resp.UserUUID,
			Login:    resp.Login,
			Email:    resp.Email,
		},
		NotificationMethods: notificationMethods,
	}
}
