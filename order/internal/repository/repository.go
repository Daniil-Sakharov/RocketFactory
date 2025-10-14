package repository

import (
	"context"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	Get(ctx context.Context, orderUUID string) (*domain.Order, error)
	Update(ctx context.Context, order *domain.Order) error
}
