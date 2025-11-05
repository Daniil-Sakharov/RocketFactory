package model

import "errors"

var (
	ErrLoginEmpty         = errors.New("login cannot be empty")
	ErrLoginInvalidLength = errors.New("login must be between 3 and 50 characters")
	ErrLoginAlreadyExists = errors.New("login already exists")

	ErrEmailRequired      = errors.New("email is required")
	ErrEmailSameAsCurrent = errors.New("new email is the same as current")
	ErrEmailAlreadyExists = errors.New("email already exists")

	ErrPasswordRequired      = errors.New("password is required")
	ErrPasswordSameAsCurrent = errors.New("new password is the same as current")

	ErrTooManyNotificationMethods      = errors.New("maximum 5 notification methods allowed")
	ErrNotificationMethodAlreadyExists = errors.New("notification method already exists")
	ErrNotificationMethodNotFound      = errors.New("notification method not found")

	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserUUIDRequired  = errors.New("user_uuid is required")

	ErrSessionNotFound     = errors.New("session not found")
	ErrSessionExpired      = errors.New("session expired")
	ErrSessionUUIDRequired = errors.New("session_uuid is required")
)
