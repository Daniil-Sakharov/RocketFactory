package entity

import (
	"time"

	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model/vo"
)

func (u *User) UpdateNotificationMethods(methods []vo.NotificationMethod) error {
	if len(methods) > 5 {
		return model.ErrTooManyNotificationMethods
	}

	u.notificationMethods = methods
	u.updatedAt = time.Now()

	return nil
}

func (u *User) AddNotificationMethod(method vo.NotificationMethod) error {
	if len(u.notificationMethods) >= 5 {
		return model.ErrTooManyNotificationMethods
	}

	for _, existing := range u.notificationMethods {
		if existing.ProviderName() == method.ProviderName() {
			return model.ErrNotificationMethodAlreadyExists
		}
	}

	u.notificationMethods = append(u.notificationMethods, method)
	u.updatedAt = time.Now()

	return nil
}

func (u *User) RemoveNotificationMethod(providerName string) error {
	newMethods := make([]vo.NotificationMethod, 0, len(u.notificationMethods))
	found := false

	for _, method := range u.notificationMethods {
		if method.ProviderName() != providerName {
			newMethods = append(newMethods, method)
		} else {
			found = true
		}
	}

	if !found {
		return model.ErrNotificationMethodNotFound
	}

	u.notificationMethods = newMethods
	u.updatedAt = time.Now()

	return nil
}

func (u *User) HasNotificationMethod(providerName string) bool {
	for _, method := range u.notificationMethods {
		if method.ProviderName() == providerName {
			return true
		}
	}
	return false
}
