package grpc

import (
	"context"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/domain"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/service/dto"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter *domain.PartsFilter) ([]*domain.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, req *dto.PayOrderClientRequest) (*dto.PayOrderClientResponse, error)
}
