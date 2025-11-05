package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model/entity"
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/service/dto"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
)

func (s *service) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	if req.Login == "" {
		logger.Warn(ctx, "login attempt with empty login")
		return nil, fmt.Errorf("login is required")
	}

	if req.Password == "" {
		logger.Warn(ctx, "login attempt with empty password")
		return nil, fmt.Errorf("password is required")
	}

	user, err := s.userRepository.GetByLogin(ctx, req.Login)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			logger.Warn(ctx, "login attempt with non-existent login",
				zap.String("login", req.Login),
			)
			return nil, model.ErrUserNotFound
		}
		logger.Error(ctx, "failed to get user by login",
			zap.String("login", req.Login),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if err = user.Password().CompareWith(req.Password); err != nil {
		logger.Warn(ctx, "login attempt with invalid password",
			zap.String("login", req.Login),
			zap.String("user_uuid", user.UserUUID()),
		)
		return nil, fmt.Errorf("invalid credentials")
	}

	sessionUUID := uuid.New().String()
	session := entity.NewSession(sessionUUID, user.UserUUID(), s.sessionTTL)

	if err = s.sessionRepository.Create(ctx, session, s.sessionTTL); err != nil {
		logger.Error(ctx, "failed to create session in redis",
			zap.String("user_uuid", user.UserUUID()),
			zap.String("session_uuid", sessionUUID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	logger.Info(ctx, "user logged in successfully",
		zap.String("user_uuid", user.UserUUID()),
		zap.String("login", user.Login()),
		zap.String("session_uuid", sessionUUID),
	)

	return &dto.LoginResponse{
		SessionUUID: sessionUUID,
	}, nil
}
