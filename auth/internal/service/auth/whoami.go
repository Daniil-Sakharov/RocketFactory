package auth

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/service/dto"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
)

func (s *service) Whoami(ctx context.Context, req *dto.WhoamiRequest) (*dto.WhoamiResponse, error) {
	if req.SessionUUID == "" {
		logger.Warn(ctx, "whoami attempt with empty session_uuid")
		return nil, model.ErrSessionUUIDRequired
	}

	session, err := s.sessionRepository.Get(ctx, req.SessionUUID)
	if err != nil {
		if errors.Is(err, model.ErrSessionNotFound) {
			logger.Warn(ctx, "session not found",
				zap.String("session_uuid", req.SessionUUID),
			)
			return nil, model.ErrSessionNotFound
		}

		logger.Error(ctx, "failed to get session from redis",
			zap.String("session_uuid", req.SessionUUID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if session.IsExpired() {
		logger.Warn(ctx, "session expired",
			zap.String("session_uuid", req.SessionUUID),
			zap.Time("expires_at", session.ExpiresAt()),
		)
		return nil, model.ErrSessionNotFound
	}

	user, err := s.userRepository.Get(ctx, session.UserUUID())
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			logger.Error(ctx, "session exists but user not found (data inconsistency)",
				zap.String("session_uuid", req.SessionUUID),
				zap.String("user_uuid", session.UserUUID()),
			)
			return nil, model.ErrUserNotFound
		}

		logger.Error(ctx, "failed to get user from database",
			zap.String("user_uuid", session.UserUUID()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	logger.Debug(ctx, "whoami successful",
		zap.String("session_uuid", req.SessionUUID),
		zap.String("user_uuid", user.UserUUID()),
		zap.String("login", user.Login()),
	)

	return &dto.WhoamiResponse{
		UserUUID: user.UserUUID(),
		Login:    user.Login(),
		Email:    user.Email().Value(),
	}, nil
}
