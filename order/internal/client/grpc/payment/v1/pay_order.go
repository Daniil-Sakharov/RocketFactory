package v1

import (
	"context"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/client/converter"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/dto"
	generatedPayment "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, req *dto.PayOrderClientRequest) (*dto.PayOrderClientResponse, error) {
	response, err := c.generatedClient.PayOrder(ctx, &generatedPayment.PayOrderRequest{
		OrderUuid:     req.OrderUUID,
		UserUuid:      req.UserUUID,
		PaymentMethod: converter.PaymentMethodToProto(req.PaymentMethod),
	})
	if err != nil {
		return nil, err
	}
	return converter.PaymentResponseFromProto(response), nil
}
