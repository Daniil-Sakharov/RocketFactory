package v1

import (
	"context"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/converter"
	"github.com/Daniil-Sakharov/RocketFactory/order/pkg/apierrors"
	orderV1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	serviceReq := converter.PayOrderRequestToServiceModel(*req, params.OrderUUID.String())

	order, err := a.service.Pay(ctx, serviceReq)
	if err != nil {
		return apierrors.MapToPayOrderError(err), nil
	}

	return converter.PayOrderResponseFromEntity(order), nil
}
