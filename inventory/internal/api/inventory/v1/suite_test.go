package v1

import (
	"context"
	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/service/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServiceSuite struct {
	suite.Suite
	ctx            context.Context
	partService *mocks.PartService
	api *api
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.partService = mocks.NewPartService(s.T())

	s.api = NewAPI(
		s.partService,
	)
}

func (s *ServiceSuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
