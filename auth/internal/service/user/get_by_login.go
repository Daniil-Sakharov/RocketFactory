package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/service/dto"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
)

func (s *service) GetByLogin(ctx context.Context, login string) (*dto.GetUserResponse, error) {
	if login == "" {
		logger.Warn(ctx, "get user by login attempt with empty login")
		return nil, model.ErrLoginEmpty
	}

	user, err := s.userRepository.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			logger.Warn(ctx, "user not found by login",
				zap.String("login", login),
			)
			return nil, model.ErrUserNotFound
		}

		logger.Error(ctx, "failed to get user by login from database",
			zap.String("login", login),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get user by login: %w", err)
	}

	return &dto.GetUserResponse{
		UserUUID:            user.UserUUID(),
		Login:               user.Login(),
		Email:               user.Email().Value(),
		NotificationMethods: notificationMethodsVOToDTO(user.NotificationMethods()),
		CreatedAt:           user.CreatedAt().Format(time.RFC3339),
		UpdatedAt:           user.UpdatedAt().Format(time.RFC3339),
	}, nil
}
