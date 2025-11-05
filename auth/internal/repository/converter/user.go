package converter

import (
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model/entity"
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model/vo"
	repoModel "github.com/Daniil-Sakharov/RocketFactory/auth/internal/repository/model"
)

func EntityToRepositoryUser(user *entity.User) *repoModel.User {
	notificationMethods := make(repoModel.NotificationMethodsDB, 0, len(user.NotificationMethods()))
	for _, method := range user.NotificationMethods() {
		notificationMethods = append(notificationMethods, repoModel.NotificationMethodDB{
			ProviderName: method.ProviderName(),
			Target:       method.Target(),
		})
	}

	return &repoModel.User{
		UserUUID:            user.UserUUID(),
		Login:               user.Login(),
		PasswordHash:        user.Password().Hash(),
		Email:               user.Email().Value(),
		NotificationMethods: notificationMethods,
		CreatedAt:           user.CreatedAt(),
		UpdatedAt:           user.UpdatedAt(),
	}
}

func RepositoryToEntityUser(repoUser *repoModel.User) *entity.User {
	//nolint:gosec // Data from DB is already validated
	email, _ := vo.NewEmail(repoUser.Email)
	//nolint:gosec // Data from DB is already validated
	password, _ := vo.NewPasswordFromHash(repoUser.PasswordHash)

	notificationMethods := make([]vo.NotificationMethod, 0, len(repoUser.NotificationMethods))
	for _, methodDB := range repoUser.NotificationMethods {
		//nolint:gosec // Data from DB is already validated
		method, _ := vo.NewNotificationMethod(methodDB.ProviderName, methodDB.Target)
		notificationMethods = append(notificationMethods, *method)
	}

	return entity.RestoreUser(
		repoUser.UserUUID,
		repoUser.Login,
		password,
		email,
		notificationMethods,
		repoUser.CreatedAt,
		repoUser.UpdatedAt,
	)
}
