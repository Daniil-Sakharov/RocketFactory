package v1

import (
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/service"
	orderV1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/openapi/order/v1"
)

var _ orderV1.Handler = (*api)(nil)

type api struct {
	service service.OrderService
}

func NewAPI(service service.OrderService) *api {
	return &api{
		service: service,
	}
}
