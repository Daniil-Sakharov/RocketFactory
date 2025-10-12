package order

import (
	"context"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/entity"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/repository/converter"
)

func (r *repository) Get(_ context.Context, orderUUID string) (*entity.Order, error) {
	r.mu.RLock()
	repoOrder, exist := r.repo[orderUUID]
	if !exist {
		r.mu.RUnlock()
		return nil, model.ErrOrderNotFound
	}
	r.mu.RUnlock()
	order := converter.OrderFromRepoModel(repoOrder)
	return order, nil
}
