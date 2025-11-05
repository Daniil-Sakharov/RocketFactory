package v1

import (
	"context"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/converter"
	api2 "github.com/Daniil-Sakharov/RocketFactory/order/internal/converter/api"
	orderV1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	uuid := converter.GetOrderRequestToServiceModel(params.OrderUUID.String())
	order, err := a.service.Get(ctx, uuid)
	if err != nil {
		return api2.MapToGetOrderError(err), nil
	}

	return converter.GetOrderResponseFromEntity(order), nil
}
