package part

import (
	"context"
	"errors"
	"fmt"

	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/model"
)

func (s *service) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	part, err := s.partRepository.GetPart(ctx, uuid)
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, model.ErrPartNotFound
		}
		return nil, fmt.Errorf("failed to get part")
	}

	return part, nil
}
