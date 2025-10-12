package v1

import (
	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/service"
	inventoryv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryv1.UnimplementedInventoryServiceServer
	partService service.PartService
}

func NewAPI(partService service.PartService) *api {
	return &api{
		partService: partService,
	}
}
