package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/domain"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/dto"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/vo"
)

func (s *service) Pay(ctx context.Context, req *dto.PayOrderRequest) (*domain.Order, error) {
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
		return nil, fmt.Errorf("failed to access pay")
	}

	newOrder := &domain.Order{
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
		if errors.Is(err, model.ErrOrderAlreadyPaid) {
			return nil, err
		}
		if errors.Is(err, model.ErrOrderAlreadyCancelled) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	produceOrder := &domain.OrderProduceEvent{
		EventUUID:       uuid.New().String(),
		OrderUUID:       newOrder.OrderUUID,
		UserUUID:        newOrder.UserUUID,
		PaymentMethod:   string(newOrder.PaymentMethod),
		TransactionUUID: newOrder.TransactionUUID,
	}

	err = s.orderProducer.PublishOrderPaid(ctx, produceOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to produce order: %w", err)
	}

	return newOrder, nil
}
