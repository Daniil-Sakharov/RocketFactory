package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	clientMocks "github.com/Daniil-Sakharov/RocketFactory/order/internal/client/grpc/mocks"
	repoMocks "github.com/Daniil-Sakharov/RocketFactory/order/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite
	ctx             context.Context
	orderRepository *repoMocks.OrderRepository
	inventoryClient *clientMocks.InventoryClient
	paymentClient   *clientMocks.PaymentClient
	service         *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.orderRepository = repoMocks.NewOrderRepository(s.T())
	s.inventoryClient = clientMocks.NewInventoryClient(s.T())
	s.paymentClient = clientMocks.NewPaymentClient(s.T())

	s.service = NewService(
		s.orderRepository,
		s.inventoryClient,
		s.paymentClient,
	)
}

func (s *ServiceSuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
