package user

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model/entity"
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model/vo"
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/service/dto"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
)

func (s *service) Update(ctx context.Context, req *dto.UpdateUserRequest) error {
	if err := s.validateUpdateRequest(ctx, req); err != nil {
		return err
	}

	user, err := s.getUserForUpdate(ctx, req.UserUUID)
	if err != nil {
		return err
	}

	needsSessionInvalidation := false

	if req.Email != nil {
		changed, err := s.updateUserEmail(ctx, user, req)
		if err != nil {
			return err
		}
		needsSessionInvalidation = needsSessionInvalidation || changed
	}

	if req.Password != nil {
		changed, err := s.updateUserPassword(ctx, user, req)
		if err != nil {
			return err
		}
		needsSessionInvalidation = needsSessionInvalidation || changed
	}

	if req.NotificationMethods != nil {
		if err := s.updateUserNotificationMethods(ctx, user, req); err != nil {
			return err
		}
	}

	if err := s.saveUserUpdates(ctx, user, req.UserUUID); err != nil {
		return err
	}

	s.invalidateSessionsIfNeeded(ctx, user.UserUUID(), needsSessionInvalidation)

	logger.Info(ctx, "user updated successfully",
		zap.String("user_uuid", req.UserUUID),
		zap.Bool("email_changed", req.Email != nil && needsSessionInvalidation),
		zap.Bool("password_changed", req.Password != nil && needsSessionInvalidation),
		zap.Bool("notifications_changed", req.NotificationMethods != nil),
	)

	return nil
}

func (s *service) validateUpdateRequest(ctx context.Context, req *dto.UpdateUserRequest) error {
	if req.UserUUID == "" {
		logger.Warn(ctx, "update user attempt with empty user_uuid")
		return model.ErrUserUUIDRequired
	}

	if req.Email == nil && req.Password == nil && req.NotificationMethods == nil {
		logger.Warn(ctx, "update user attempt with no fields to update",
			zap.String("user_uuid", req.UserUUID),
		)
		return fmt.Errorf("nothing to update")
	}

	return nil
}

func (s *service) getUserForUpdate(ctx context.Context, userUUID string) (*entity.User, error) {
	user, err := s.userRepository.Get(ctx, userUUID)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			logger.Warn(ctx, "user not found for update", zap.String("user_uuid", userUUID))
			return nil, model.ErrUserNotFound
		}
		logger.Error(ctx, "failed to get user for update", zap.String("user_uuid", userUUID), zap.Error(err))
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (s *service) updateUserEmail(ctx context.Context, user *entity.User, req *dto.UpdateUserRequest) (bool, error) {
	newEmail, err := vo.NewEmail(*req.Email)
	if err != nil {
		logger.Warn(ctx, "invalid email format for update",
			zap.String("user_uuid", req.UserUUID),
			zap.String("email", *req.Email),
			zap.Error(err),
		)
		return false, fmt.Errorf("invalid email: %w", err)
	}

	if err := user.UpdateEmail(newEmail); err != nil {
		if errors.Is(err, model.ErrEmailSameAsCurrent) {
			logger.Debug(ctx, "email not changed (same as current)", zap.String("user_uuid", req.UserUUID))
			return false, nil
		}
		logger.Error(ctx, "failed to update email", zap.String("user_uuid", req.UserUUID), zap.Error(err))
		return false, fmt.Errorf("failed to update email: %w", err)
	}

	logger.Info(ctx, "email updated successfully",
		zap.String("user_uuid", req.UserUUID),
		zap.String("new_email", *req.Email),
	)
	return true, nil
}

func (s *service) updateUserPassword(ctx context.Context, user *entity.User, req *dto.UpdateUserRequest) (bool, error) {
	newPassword, err := vo.NewPasswordFromPlaintext(*req.Password)
	if err != nil {
		logger.Warn(ctx, "invalid password for update", zap.String("user_uuid", req.UserUUID), zap.Error(err))
		return false, fmt.Errorf("invalid password: %w", err)
	}

	if err := user.UpdatePassword(newPassword); err != nil {
		if errors.Is(err, model.ErrPasswordSameAsCurrent) {
			logger.Debug(ctx, "password not changed (same as current)", zap.String("user_uuid", req.UserUUID))
			return false, nil
		}
		logger.Error(ctx, "failed to update password", zap.String("user_uuid", req.UserUUID), zap.Error(err))
		return false, fmt.Errorf("failed to update password: %w", err)
	}

	logger.Info(ctx, "password updated successfully", zap.String("user_uuid", req.UserUUID))
	return true, nil
}

func (s *service) updateUserNotificationMethods(ctx context.Context, user *entity.User, req *dto.UpdateUserRequest) error {
	notificationMethods, err := notificationMethodsDTOToVO(req.NotificationMethods)
	if err != nil {
		logger.Warn(ctx, "invalid notification methods for update",
			zap.String("user_uuid", req.UserUUID),
			zap.Error(err),
		)
		return fmt.Errorf("invalid notification methods: %w", err)
	}

	if err := user.UpdateNotificationMethods(notificationMethods); err != nil {
		logger.Error(ctx, "failed to update notification methods",
			zap.String("user_uuid", req.UserUUID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update notification methods: %w", err)
	}

	logger.Info(ctx, "notification methods updated successfully",
		zap.String("user_uuid", req.UserUUID),
		zap.Int("methods_count", len(notificationMethods)),
	)
	return nil
}

func (s *service) saveUserUpdates(ctx context.Context, user *entity.User, userUUID string) error {
	if err := s.userRepository.Update(ctx, user); err != nil {
		if errors.Is(err, model.ErrEmailAlreadyExists) {
			logger.Warn(ctx, "email already exists", zap.String("user_uuid", userUUID))
			return model.ErrEmailAlreadyExists
		}
		logger.Error(ctx, "failed to update user in database", zap.String("user_uuid", userUUID), zap.Error(err))
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (s *service) invalidateSessionsIfNeeded(ctx context.Context, userUUID string, needsInvalidation bool) {
	if !needsInvalidation {
		return
	}

	logger.Info(ctx, "invalidating all user sessions due to email/password change", zap.String("user_uuid", userUUID))

	if err := s.sessionRepository.DeleteAllByUserUUID(ctx, userUUID); err != nil {
		logger.Warn(ctx, "failed to invalidate sessions after email/password change",
			zap.String("user_uuid", userUUID),
			zap.Error(err),
		)
		return
	}

	logger.Info(ctx, "all user sessions invalidated successfully", zap.String("user_uuid", userUUID))
}
