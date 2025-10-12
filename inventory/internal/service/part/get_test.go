package part

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/model"
)

func (s *ServiceSuite) TestGetPartSuccess() {
	var (
		partUUID = gofakeit.UUID()

		name        = "Advanced Navigation System"
		description = "High-precision navigation system with GPS and inertial guidance"
		price       = 5000000.00
		stock       = int64(12)
		category    = model.CATEGORY_ENGINE

		dimensions = &model.Dimensions{
			Length: 50.0,
			Width:  40.0,
			Height: 30.0,
			Weight: 25.5,
		}

		manufacturer = &model.Manufacturer{
			Name:    "Honeywell Aerospace",
			Country: "USA",
			Website: "https://www.honeywell.com",
		}

		tags = []string{"navigation", "gps", "precision", "electronics"}

		metadata = map[string]interface{}{
			"accuracy":       "0.01m",
			"update_rate":    "100Hz",
			"power_draw":     "50W",
			"temperature":    "-40°C to +85°C",
			"certifications": []string{"AS9100", "ISO 9001"},
		}

		createdAt = time.Now().Add(-120 * 24 * time.Hour)
		updatedAt = time.Now().Add(-3 * 24 * time.Hour)

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

	s.partRepository.On("GetPart", s.ctx, partUUID).Return(expectedPart, nil)

	part, err := s.service.GetPart(s.ctx, partUUID)
	s.Require().NoError(err)
	s.Require().Equal(expectedPart, part)
}

func (s *ServiceSuite) TestGetPartNotFound() {
	var (
		partUUID = gofakeit.UUID()
	)

	s.partRepository.On("GetPart", s.ctx, partUUID).Return(nil, model.ErrPartNotFound)

	part, err := s.service.GetPart(s.ctx, partUUID)
	s.Require().Error(err)
	s.Require().Nil(part)
	s.Require().ErrorIs(err, model.ErrPartNotFound)
}
