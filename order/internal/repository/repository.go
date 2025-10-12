package repository

import (
	"context"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/entity"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error
	Get(ctx context.Context, orderUUID string) (*entity.Order, error)
	Update(ctx context.Context, order *entity.Order) error
}
