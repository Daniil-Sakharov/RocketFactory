package order

import (
	"context"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/entity"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/repository/converter"
)

func (r *repository) Update(_ context.Context, order *entity.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.repo[order.OrderUUID]
	if !ok {
		return model.ErrOrderNotFound
	}

	repoOrder := converter.OrderToRepoModel(order)
	r.repo[order.OrderUUID] = repoOrder

	return nil
}
