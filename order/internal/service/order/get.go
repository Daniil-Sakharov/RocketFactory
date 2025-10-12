package order

import (
	"context"
	"errors"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/dto"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/entity"
)

func (s *service) Get(ctx context.Context, req *dto.GetOrderRequest) (*entity.Order, error) {
	order, err := s.orderRepository.Get(ctx, req.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return nil, err
		}
		return nil, model.ErrUnknownError
	}
	return order, nil
}
