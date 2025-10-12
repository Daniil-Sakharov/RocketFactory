package part

import (
	"context"

	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	part, err := s.partRepository.ListParts(ctx, filter)
	if err != nil {
		return nil, err
	}
	return part, nil
}
