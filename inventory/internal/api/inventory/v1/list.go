package v1

import (
	"context"

	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/converter"
	inventoryv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	filter := converter.FilterFromProto(req.Filter)

	parts, err := a.partService.ListParts(ctx, filter)
	if err != nil {
		return nil, err
	}

	protoParts := converter.PartsToProto(parts)

	return &inventoryv1.ListPartsResponse{
		Parts: protoParts,
	}, nil
}
