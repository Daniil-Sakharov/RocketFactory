package order

import (
	"context"
	"errors"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/dto"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/vo"
)

func (s *service) Cancel(ctx context.Context, req *dto.CancelOrderRequest) error {
	order, err := s.orderRepository.Get(ctx, req.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return err
		}
		return model.ErrUnknownError
	}
	if order.Status == vo.OrderStatusPAID {
		return model.ErrOrderAlreadyPaid
	}
	order.Status = vo.OrderStatusCANCELLED
	err = s.orderRepository.Update(ctx, order)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return err
		}
		return model.ErrUnknownError
	}
	return nil
}
