package order

import (
	"context"
	"errors"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/dto"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/entity"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/vo"
)

func (s *service) Pay(ctx context.Context, req *dto.PayOrderRequest) (*entity.Order, error) {
	order, err := s.orderRepository.Get(ctx, req.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return nil, err
		}
		return nil, model.ErrUnknownError
	}
	response, err := s.paymentClient.PayOrder(ctx, &dto.PayOrderClientRequest{
		OrderUUID:     order.OrderUUID,
		UserUUID:      order.UserUUID,
		PaymentMethod: req.PaymentMethod,
	})
	if err != nil {
		return nil, errors.New("failed to access pay")
	}

	newOrder := &entity.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: response.TransactionUUID,
		PaymentMethod:   req.PaymentMethod,
		Status:          vo.OrderStatusPAID,
	}

	err = s.orderRepository.Update(ctx, newOrder)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return nil, err
		}
		return nil, model.ErrUnknownError
	}

	return newOrder, nil
}
