package v1

import (
	"context"

	clientConverter "github.com/Daniil-Sakharov/RocketFactory/order/internal/client/converter"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/entity"
	generatedInventoryV1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter *entity.PartsFilter) ([]*entity.Part, error) {
	response, err := c.generatedClient.ListParts(ctx, &generatedInventoryV1.ListPartsRequest{
		Filter: clientConverter.FilterToProto(filter),
	})
	if err != nil {
		return nil, err
	}
	return clientConverter.PartsFromProto(response.Parts), nil
}
