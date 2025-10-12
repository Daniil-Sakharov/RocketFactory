package v1

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/model"
	inventoryv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/inventory/v1"
)

func (s *ServiceSuite) TestGetPartSuccess() {
	var (
		partUUID = gofakeit.UUID()

		request = &inventoryv1.GetPartRequest{
			Uuid: partUUID,
		}

		name        = "Flight Control System"
		description = "Advanced flight control system with redundancy and autopilot"
		price       = 12000000.00
		stock       = int64(8)
		category    = model.CATEGORY_ENGINE

		dimensions = &model.Dimensions{
			Length: 120.0,
			Width:  80.0,
			Height: 60.0,
			Weight: 150.0,
		}

		manufacturer = &model.Manufacturer{
			Name:    "Northrop Grumman",
			Country: "USA",
			Website: "https://www.northropgrumman.com",
		}

		tags = []string{"flight", "control", "autopilot", "redundancy"}

		metadata = map[string]interface{}{
			"redundancy_level": "triple",
			"processor":        "RAD750",
			"operating_system": "VxWorks",
			"mtbf_hours":       100000,
			"certifications":   []string{"DO-178C", "MIL-STD-1553"},
		}

		createdAt = time.Now().Add(-180 * 24 * time.Hour)
		updatedAt = time.Now().Add(-7 * 24 * time.Hour)

		expectedPart = &model.Part{
			Uuid:          partUUID,
			Name:          name,
			Description:   description,
			Price:         price,
			StockQuantity: stock,
			Category:      category,
			Dimensions:    dimensions,
			Manufacturer:  manufacturer,
			Tags:          tags,
			Metadata:      metadata,
			CreatedAt:     &createdAt,
			UpdatedAt:     &updatedAt,
		}
	)

	s.partService.On("GetPart", s.ctx, partUUID).Return(expectedPart, nil)

	response, err := s.api.GetPart(s.ctx, request)
	s.Require().NoError(err)
	s.Require().NotNil(response)
	s.Require().NotNil(response.Part)
	s.Require().Equal(partUUID, response.Part.Uuid)
	s.Require().Equal(name, response.Part.Name)
	s.Require().Equal(description, response.Part.Description)
}

func (s *ServiceSuite) TestGetPartNotFound() {
	var (
		partUUID = gofakeit.UUID()

		request = &inventoryv1.GetPartRequest{
			Uuid: partUUID,
		}
	)

	s.partService.On("GetPart", s.ctx, partUUID).Return(nil, model.ErrPartNotFound)

	response, err := s.api.GetPart(s.ctx, request)
	s.Require().Error(err)
	s.Require().Nil(response)
}
