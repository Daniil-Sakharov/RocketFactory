package order

import (
	def "github.com/Daniil-Sakharov/RocketFactory/order/internal/repository"
	repoModel "github.com/Daniil-Sakharov/RocketFactory/order/internal/repository/model"
	"sync"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	repo map[string]*repoModel.Order
}

func NewRepository() *repository {
	return &repository{
		repo: make(map[string]*repoModel.Order),
	}
}
