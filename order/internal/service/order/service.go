package order

import (
	client "github.com/Daniil-Sakharov/RocketFactory/order/internal/client/grpc"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/repository"
	def "github.com/Daniil-Sakharov/RocketFactory/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepository repository.OrderRepository
	inventoryClient client.InventoryClient
	paymentClient   client.PaymentClient
	orderProducer   def.OrderProducerService
}

func NewService(
	orderRepository repository.OrderRepository,
	inventoryClient client.InventoryClient,
	paymentClient client.PaymentClient,
	orderProducer def.OrderProducerService,
) *service {
	return &service{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		orderProducer:   orderProducer,
	}
}
