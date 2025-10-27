package service

import (
	"context"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/domain"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/dto"
)

type OrderService interface {
	Create(ctx context.Context, req *dto.CreateOrderRequest) (*domain.Order, error)
	Pay(ctx context.Context, req *dto.PayOrderRequest) (*domain.Order, error)
	Get(ctx context.Context, req *dto.GetOrderRequest) (*domain.Order, error)
	Cancel(ctx context.Context, req *dto.CancelOrderRequest) error
}

type AssemblyConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderProducerService interface {
	PublishOrderPaid(ctx context.Context, event *domain.OrderProduceEvent) error
}
