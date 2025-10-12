package v1

import (
	def "github.com/Daniil-Sakharov/RocketFactory/order/internal/client/grpc"
	generatedInventoryV1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/inventory/v1"
)

var _ def.InventoryClient = (*client)(nil)

type client struct {
	generatedClient generatedInventoryV1.InventoryServiceClient
}

func NewClient(generatedClient generatedInventoryV1.InventoryServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}
