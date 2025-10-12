package part

import (
	"context"

	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/repository/converter"
)

func (r *repository) GetPart(_ context.Context, uuid string) (*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoPart, ok := r.data[uuid]
	if !ok {
		return nil, model.ErrPartNotFound
	}

	return converter.PartToModel(&repoPart), nil
}
