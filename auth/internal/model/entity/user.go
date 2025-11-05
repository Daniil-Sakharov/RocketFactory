package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model/vo"
)

type User struct {
	userUUID            string
	login               string
	password            *vo.Password
	email               *vo.Email
	notificationMethods []vo.NotificationMethod
	createdAt           time.Time
	updatedAt           time.Time
}

func NewUser(
	login string,
	password *vo.Password,
	email *vo.Email,
	notificationMethods []vo.NotificationMethod,
) (*User, error) {
	if login == "" {
		return nil, errors.New("login cannot be empty")
	}

	if len(login) < 3 || len(login) > 50 {
		return nil, errors.New("login must be between 3 and 50 characters")
	}

	if len(notificationMethods) > 5 {
		return nil, errors.New("maximum 5 notification methods allowed")
	}

	return &User{
		userUUID:            uuid.New().String(),
		login:               login,
		email:               email,
		password:            password,
		notificationMethods: notificationMethods,
		createdAt:           time.Now(),
		updatedAt:           time.Now(),
	}, nil
}

func (u *User) UserUUID() string {
	return u.userUUID
}

func (u *User) Login() string {
	return u.login
}

func (u *User) Password() *vo.Password {
	return u.password
}

func (u *User) Email() *vo.Email {
	return u.email
}

func (u *User) NotificationMethods() []vo.NotificationMethod {
	methods := make([]vo.NotificationMethod, len(u.notificationMethods))
	copy(methods, u.notificationMethods)
	return methods
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func RestoreUser(
	userUUID string,
	login string,
	password *vo.Password,
	email *vo.Email,
	notificationMethods []vo.NotificationMethod,
	createdAt time.Time,
	updatedAt time.Time,
) *User {
	return &User{
		userUUID:            userUUID,
		login:               login,
		password:            password,
		email:               email,
		notificationMethods: notificationMethods,
		createdAt:           createdAt,
		updatedAt:           updatedAt,
	}
}
