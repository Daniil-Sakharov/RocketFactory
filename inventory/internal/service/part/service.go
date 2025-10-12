package part

import (
	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/repository"
	def "github.com/Daniil-Sakharov/RocketFactory/inventory/internal/service"
)

var _ def.PartService = (*service)(nil)

type service struct {
	partRepository repository.PartRepository
}

func NewService(partRepository repository.PartRepository) *service {
	return &service{
		partRepository: partRepository,
	}
}
