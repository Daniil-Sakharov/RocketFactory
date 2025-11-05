package converter

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model/vo"
)

func MapErrorToGRPC(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, model.ErrUserNotFound) {
		return status.Error(codes.NotFound, "user not found")
	}

	if errors.Is(err, model.ErrSessionNotFound) {
		return status.Error(codes.Unauthenticated, "session not found or expired")
	}

	if errors.Is(err, model.ErrLoginAlreadyExists) {
		return status.Error(codes.AlreadyExists, "login already exists")
	}

	if errors.Is(err, model.ErrEmailAlreadyExists) {
		return status.Error(codes.AlreadyExists, "email already exists")
	}

	if errors.Is(err, model.ErrLoginEmpty) {
		return status.Error(codes.InvalidArgument, "login is required")
	}

	if errors.Is(err, model.ErrPasswordRequired) {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if errors.Is(err, model.ErrEmailRequired) {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if errors.Is(err, model.ErrUserUUIDRequired) {
		return status.Error(codes.InvalidArgument, "user_uuid is required")
	}

	if errors.Is(err, model.ErrSessionUUIDRequired) {
		return status.Error(codes.InvalidArgument, "session_uuid is required")
	}

	if errors.Is(err, model.ErrLoginInvalidLength) {
		return status.Error(codes.InvalidArgument, "login must be between 3 and 50 characters")
	}

	if errors.Is(err, vo.ErrEmailInvalidFormat) {
		return status.Error(codes.InvalidArgument, "invalid email format")
	}

	if errors.Is(err, vo.ErrPasswordTooShort) {
		return status.Error(codes.InvalidArgument, "password must be at least 8 characters")
	}

	if errors.Is(err, vo.ErrPasswordEmpty) {
		return status.Error(codes.InvalidArgument, "password cannot be empty")
	}

	if errors.Is(err, vo.ErrNotificationProviderInvalid) {
		return status.Error(codes.InvalidArgument, "invalid notification provider")
	}

	if errors.Is(err, vo.ErrNotificationTargetInvalid) {
		return status.Error(codes.InvalidArgument, "invalid notification target format")
	}

	if errors.Is(err, model.ErrTooManyNotificationMethods) {
		return status.Error(codes.InvalidArgument, "maximum 5 notification methods allowed")
	}

	if errors.Is(err, model.ErrEmailSameAsCurrent) {
		return status.Error(codes.FailedPrecondition, "new email is the same as current")
	}

	if errors.Is(err, model.ErrPasswordSameAsCurrent) {
		return status.Error(codes.FailedPrecondition, "new password is the same as current")
	}

	return status.Error(codes.Internal, "internal server error")
}
