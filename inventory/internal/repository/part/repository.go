package part

import (
	"sync"

	def "github.com/Daniil-Sakharov/RocketFactory/inventory/internal/repository"
	repoModel "github.com/Daniil-Sakharov/RocketFactory/inventory/internal/repository/model"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]repoModel.Part
}

func NewRepository() *repository {
	return &repository{
		data: make(map[string]repoModel.Part),
	}
}
