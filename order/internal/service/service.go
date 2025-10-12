package service

import (
	"context"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/dto"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/entity"
)

type OrderService interface {
	Create(ctx context.Context, req *dto.CreateOrderRequest) (*entity.Order, error)
	Pay(ctx context.Context, req *dto.PayOrderRequest) (*entity.Order, error)
	Get(ctx context.Context, req *dto.GetOrderRequest) (*entity.Order, error)
	Cancel(ctx context.Context, req *dto.CancelOrderRequest) error
}
