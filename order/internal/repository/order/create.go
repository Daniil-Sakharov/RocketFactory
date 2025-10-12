package order

import (
	"context"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/entity"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/repository/converter"
)

func (r *repository) Create(_ context.Context, order *entity.Order) error {
	repoOrder := converter.OrderToRepoModel(order)
	r.mu.Lock()
	if _, exist := r.repo[repoOrder.OrderUUID]; exist {
		return model.ErrOrderAlreadyExist
	}
	r.repo[repoOrder.OrderUUID] = repoOrder
	r.mu.Unlock()
	return nil
}
