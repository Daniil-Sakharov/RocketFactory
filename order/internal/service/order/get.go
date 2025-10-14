package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/domain"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/dto"
)

func (s *service) Get(ctx context.Context, req *dto.GetOrderRequest) (*domain.Order, error) {
	order, err := s.orderRepository.Get(ctx, req.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get order %w", err)
	}
	return order, nil
}
