// auth/internal/service/user/converter.go

package user

import (
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model/vo"
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/service/dto"
)

func notificationMethodsVOToDTO(vos []vo.NotificationMethod) []dto.NotificationMethodDTO {
	if len(vos) == 0 {
		return []dto.NotificationMethodDTO{}
	}

	dtos := make([]dto.NotificationMethodDTO, 0, len(vos))
	for _, vo := range vos {
		dtos = append(dtos, dto.NotificationMethodDTO{
			ProviderName: vo.ProviderName(),
			Target:       vo.Target(),
		})
	}
	return dtos
}

func notificationMethodsDTOToVO(dtos []dto.NotificationMethodDTO) ([]vo.NotificationMethod, error) {
	vos := make([]vo.NotificationMethod, 0, len(dtos))
	for _, d := range dtos {
		v, err := vo.NewNotificationMethod(d.ProviderName, d.Target)
		if err != nil {
			return nil, err
		}
		vos = append(vos, *v)
	}
	return vos, nil
}
