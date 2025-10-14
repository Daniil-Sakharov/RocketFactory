package part

import (
	"context"
	"errors"
	"fmt"

	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	part, err := s.partRepository.ListParts(ctx, filter)
	if err != nil {
		if errors.Is(err, model.ErrPartsNotFound) {
			return nil, model.ErrPartsNotFound
		}
		return nil, fmt.Errorf("failed to get parts: %w", err)
	}
	return part, nil
}
