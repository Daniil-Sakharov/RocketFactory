package part

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/model"
)

func (s *ServiceSuite) TestListPartsSuccess() {
	var (
		filter = &model.PartsFilter{
			Categories: []model.Category{model.CATEGORY_ENGINE, model.CATEGORY_FUEL},
			Tags:       []string{"rocket", "premium"},
		}

		uuid1        = gofakeit.UUID()
		name1        = "RD-180 Rocket Engine"
		description1 = "High-performance rocket engine with thrust vector control"
		price1       = 25000000.0
		stock1       = int64(5)
		category1    = model.CATEGORY_ENGINE

		dimensions1 = &model.Dimensions{
			Length: 360.0,
			Width:  310.0,
			Height: 310.0,
			Weight: 5480.0,
		}

		manufacturer1 = &model.Manufacturer{
			Name:    "Energomash",
			Country: "Russia",
			Website: "https://www.energomash.ru",
		}

		tags1 = []string{"rocket", "premium", "high-thrust"}

		metadata1 = map[string]interface{}{
			"thrust":           "4152 kN",
			"specific_impulse": "311 s",
			"fuel_type":        "RP-1/LOX",
		}

		createdAt1 = time.Now().Add(-30 * 24 * time.Hour)
		updatedAt1 = time.Now().Add(-1 * time.Hour)

		uuid2        = gofakeit.UUID()
		name2        = "Liquid Oxygen (LOX)"
		description2 = "Cryogenic oxidizer for rocket propulsion"
		price2       = 150.00
		stock2       = int64(10000)
		category2    = model.CATEGORY_FUEL

		dimensions2 = &model.Dimensions{
			Length: 200.0,
			Width:  200.0,
			Height: 400.0,
			Weight: 1141.0,
		}

		manufacturer2 = &model.Manufacturer{
			Name:    "Air Liquide",
			Country: "France",
			Website: "https://www.airliquide.com",
		}

		tags2 = []string{"rocket", "cryogenic", "oxidizer"}

		metadata2 = map[string]interface{}{
			"temperature":   "-183°C",
			"density":       "1.141 g/cm³",
			"boiling_point": "-183°C",
		}

		createdAt2 = time.Now().Add(-60 * 24 * time.Hour)
		updatedAt2 = time.Now().Add(-2 * 24 * time.Hour)

		uuid3        = gofakeit.UUID()
		name3        = "Merlin 1D Engine"
		description3 = "Reusable rocket engine with throttle capability"
		price3       = 1000000.00
		stock3       = int64(25)
		category3    = model.CATEGORY_ENGINE

		dimensions3 = &model.Dimensions{
			Length: 290.0,
			Width:  100.0,
			Height: 100.0,
			Weight: 470.0,
		}

		manufacturer3 = &model.Manufacturer{
			Name:    "SpaceX",
			Country: "USA",
			Website: "https://www.spacex.com",
		}

		tags3 = []string{"rocket", "reusable", "falcon"}

		metadata3 = map[string]interface{}{
			"thrust":           "845 kN",
			"specific_impulse": "282 s",
			"reusable":         true,
		}

		createdAt3 = time.Now().Add(-90 * 24 * time.Hour)
		updatedAt3 = time.Now().Add(-5 * time.Hour)

		expectedParts = []*model.Part{
			{
				Uuid:          uuid1,
				Name:          name1,
				Description:   description1,
				Price:         price1,
				StockQuantity: stock1,
				Category:      category1,
				Dimensions:    dimensions1,
				Manufacturer:  manufacturer1,
				Tags:          tags1,
				Metadata:      metadata1,
				CreatedAt:     &createdAt1,
				UpdatedAt:     &updatedAt1,
			},
			{
				Uuid:          uuid2,
				Name:          name2,
				Description:   description2,
				Price:         price2,
				StockQuantity: stock2,
				Category:      category2,
				Dimensions:    dimensions2,
				Manufacturer:  manufacturer2,
				Tags:          tags2,
				Metadata:      metadata2,
				CreatedAt:     &createdAt2,
				UpdatedAt:     &updatedAt2,
			},
			{
				Uuid:          uuid3,
				Name:          name3,
				Description:   description3,
				Price:         price3,
				StockQuantity: stock3,
				Category:      category3,
				Dimensions:    dimensions3,
				Manufacturer:  manufacturer3,
				Tags:          tags3,
				Metadata:      metadata3,
				CreatedAt:     &createdAt3,
				UpdatedAt:     &updatedAt3,
			},
		}
	)

	s.partRepository.On("ListParts", s.ctx, filter).Return(expectedParts, nil)

	parts, err := s.service.ListParts(s.ctx, filter)
	s.Require().NoError(err)
	s.Require().Equal(expectedParts, parts)
}

func (s *ServiceSuite) TestListPartsRepositoryError() {
	var (
		filter = &model.PartsFilter{
			Categories: []model.Category{model.CATEGORY_ENGINE},
			Tags:       []string{"rocket"},
		}

		repositoryError = model.ErrPartsNotFound
	)

	s.partRepository.On("ListParts", s.ctx, filter).Return(nil, repositoryError)

	parts, err := s.service.ListParts(s.ctx, filter)
	s.Require().Error(err)
	s.Require().Nil(parts)
	s.Require().Equal(repositoryError, err)
}
