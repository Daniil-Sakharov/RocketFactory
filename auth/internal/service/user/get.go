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

func (s *service) Get(ctx context.Context, req *dto.GetUserRequest) (*dto.GetUserResponse, error) {
	if req.UserUUID == "" {
		logger.Warn(ctx, "get with empty user_uuid")
		return nil, model.ErrUserUUIDRequired
	}

	user, err := s.userRepository.Get(ctx, req.UserUUID)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			logger.Warn(ctx, "user not found",
				zap.String("user_uuid", req.UserUUID),
			)
			return nil, model.ErrUserNotFound
		}

		logger.Error(ctx, "failed to get user from database",
			zap.String("user_uuid", req.UserUUID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get user: %w", err)
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
