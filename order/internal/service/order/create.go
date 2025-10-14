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

func (s *service) Create(ctx context.Context, req *dto.CreateOrderRequest) (*domain.Order, error) {
	if req.UserUUID == "" {
		return nil, model.ErrEmptyUserUUID
	}
	if len(req.PartUUIDs) == 0 {
		return nil, model.ErrEmptyPartUUIDs
	}

	parts, err := s.inventoryClient.ListParts(ctx, &domain.PartsFilter{
		Uuids: req.PartUUIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get parts: %w", err)
	}

	if len(parts) == 0 {
		return nil, model.ErrPartsNotFound
	}

	totalPrice := calculateTotalPrice(parts)

	newOrder := &domain.Order{
		OrderUUID:       uuid.NewString(),
		UserUUID:        req.UserUUID,
		PartUUIDs:       req.PartUUIDs,
		TotalPrice:      totalPrice,
		TransactionUUID: "",
		PaymentMethod:   vo.PaymentMethodUNKNOWN,
		Status:          vo.OrderStatusPENDINGPAYMENT,
	}
	err = s.orderRepository.Create(ctx, newOrder)
	if err != nil {
		if errors.Is(err, model.ErrOrderAlreadyExist) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	return newOrder, nil
}

func calculateTotalPrice(parts []*domain.Part) float64 {
	var total float64
	for _, part := range parts {
		total += part.Price
	}
	return total
}
