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

func (s *service) Register(ctx context.Context, req *dto.RegisterUserRequest) (*dto.RegisterUserResponse, error) {
	if req.Login == "" {
		logger.Warn(ctx, "registration attempt with empty login")
		return nil, model.ErrLoginEmpty
	}

	if req.Password == "" {
		logger.Warn(ctx, "registration attempt with empty password")
		return nil, model.ErrPasswordRequired
	}

	if req.Email == "" {
		logger.Warn(ctx, "registration attempt with empty email")
		return nil, model.ErrEmailRequired
	}

	existingUser, err := s.userRepository.GetByLogin(ctx, req.Login)
	if err != nil && !errors.Is(err, model.ErrUserNotFound) {
		logger.Error(ctx, "failed to check login uniqueness",
			zap.String("login", req.Login),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to check login uniqueness: %w", err)
	}

	if existingUser != nil {
		logger.Warn(ctx, "register attempt with existing login",
			zap.String("login", req.Login),
		)
		return nil, model.ErrLoginAlreadyExists
	}

	password, err := vo.NewPasswordFromPlaintext(req.Password)
	if err != nil {
		logger.Warn(ctx, "registration attempt with invalid password",
			zap.String("login", req.Login),
			zap.Error(err),
		)
		return nil, fmt.Errorf("invalid password: %w", err)
	}

	email, err := vo.NewEmail(req.Email)
	if err != nil {
		logger.Warn(ctx, "registration attempt with invalid email",
			zap.String("login", req.Login),
			zap.String("email", req.Email),
			zap.Error(err),
		)
		return nil, fmt.Errorf("invalid email: %w", err)
	}

	notificationMethods := make([]vo.NotificationMethod, 0, len(req.NotificationMethods))
	for _, nmDTO := range req.NotificationMethods {
		nm, err := vo.NewNotificationMethod(nmDTO.ProviderName, nmDTO.Target)
		if err != nil {
			logger.Warn(ctx, "registration attempt with invalid notification method",
				zap.String("login", req.Login),
				zap.String("provider", nmDTO.ProviderName),
				zap.Error(err),
			)
			return nil, fmt.Errorf("invalid notification method: %w", err)
		}
		notificationMethods = append(notificationMethods, *nm)
	}

	user, err := entity.NewUser(req.Login, password, email, notificationMethods)
	if err != nil {
		logger.Warn(ctx, "failed to create user entity",
			zap.String("login", req.Login),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if err = s.userRepository.Create(ctx, user); err != nil {
		if errors.Is(err, model.ErrLoginAlreadyExists) {
			logger.Warn(ctx, "login already exists (race condition)",
				zap.String("login", req.Login),
			)
			return nil, model.ErrLoginAlreadyExists
		}

		if errors.Is(err, model.ErrEmailAlreadyExists) {
			logger.Warn(ctx, "email already exists",
				zap.String("login", req.Login),
				zap.String("email", req.Email),
			)
			return nil, model.ErrEmailAlreadyExists
		}

		logger.Error(ctx, "failed to create user in database",
			zap.String("login", req.Login),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	logger.Info(ctx, "user registered successfully",
		zap.String("user_uuid", user.UserUUID()),
		zap.String("login", user.Login()),
	)

	return &dto.RegisterUserResponse{
		UserUUID: user.UserUUID(),
	}, nil
}
