package part

import (
	"context"
	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/model"
)

func (s *service) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	part, err := s.partRepository.GetPart(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return part, nil
}
