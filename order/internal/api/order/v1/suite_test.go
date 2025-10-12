package v1

import (
	"context"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/service/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type APISuite struct{
	suite.Suite
	ctx context.Context
	orderService *mocks.OrderService
	api *api
}

func (s *APISuite)SetupTest() {
	s.ctx = context.Background()
	s.orderService = mocks.NewOrderService(s.T())
	s.api = NewAPI(s.orderService)
}

func (s *APISuite) TearDown() {}

func (s *APISuite) TestServiceIntegration(t *testing.T){
	suite.Run(t, new(APISuite))
}