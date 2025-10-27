package assembly

import (
	serv "github.com/Daniil-Sakharov/RocketFactory/assembly/internal/service"
)

var _ serv.AssemblyService = (*service)(nil)

type service struct {
	shipAssembledProducer serv.ShipAssembledProducerService
}

func NewService(shipAssembledProducer serv.ShipAssembledProducerService) *service {
	return &service{
		shipAssembledProducer: shipAssembledProducer,
	}
}
