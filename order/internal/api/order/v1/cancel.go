package v1

import (
	"context"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/converter"
	"github.com/Daniil-Sakharov/RocketFactory/order/pkg/apierrors"
	orderV1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	uuid := converter.CancelOrderRequestToServiceModel(params.OrderUUID.String())
	err := a.service.Cancel(ctx, uuid)
	if err != nil {
		return apierrors.MapToCancelOrderError(err), nil
	}

	return apierrors.MapToCancelOrderError(err), nil
}
