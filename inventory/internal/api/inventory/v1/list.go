package v1

import (
	"context"

	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/converter"
	inventoryv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	parts, err := a.partService.ListParts(ctx, converter.FilterFromProto(req.Filter))
	if err != nil {
		return nil, err
	}
	return &inventoryv1.ListPartsResponse{
		Parts: converter.PartsToProto(parts),
	}, nil
}
