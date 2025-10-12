package grpc

import (
	"context"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/dto"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/entity"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter *entity.PartsFilter) ([]*entity.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, req *dto.PayOrderClientRequest) (*dto.PayOrderClientResponse, error)
}
