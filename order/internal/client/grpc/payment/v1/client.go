package v1

import (
	def "github.com/Daniil-Sakharov/RocketFactory/order/internal/client/grpc"
	generatedpaymentV1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/payment/v1"
)

var _ def.PaymentClient = (*client)(nil)

type client struct {
	generatedClient generatedpaymentV1.PaymentServiceClient
}

func NewClient(generatedClient generatedpaymentV1.PaymentServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}
