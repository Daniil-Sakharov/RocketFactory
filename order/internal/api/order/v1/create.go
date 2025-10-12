package v1

import (
	"context"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/converter"
	"github.com/Daniil-Sakharov/RocketFactory/order/pkg/apierrors"
	orderV1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	domainReq := converter.CreateOrderRequestToServiceModel(*req)

	order, err := a.service.Create(ctx, domainReq)
	if err != nil {
		return apierrors.MapToCreateOrderError(err), nil
	}

	return converter.CreateOrderResponseFromEntity(order), nil
}
