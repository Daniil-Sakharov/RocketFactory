package part

import (
	"context"
	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServiceSuite struct {
	suite.Suite
	ctx            context.Context
	partRepository *mocks.PartRepository
	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.partRepository = mocks.NewPartRepository(s.T())

	s.service = NewService(
		s.partRepository,
	)
}

func (s *ServiceSuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
