package order

import (
	"errors"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/dto"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/entity"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/vo"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuite) TestCreateOrderSuccess() {
	var (
		userUUID  = gofakeit.UUID()
		partUUID1 = gofakeit.UUID()
		partUUID2 = gofakeit.UUID()

		request = &dto.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUUIDs: []string{partUUID1, partUUID2},
		}

		filter = &entity.PartsFilter{
			Uuids: []string{partUUID1, partUUID2},
		}

		partsFromInventory = []*entity.Part{
			{
				Uuid:          partUUID1,
				Name:          "RD-180 Engine",
				Description:   "Rocket engine",
				Price:         25000000.00,
				StockQuantity: 5,
				Category:      entity.CATEGORY_ENGINE,
			},
			{
				Uuid:          partUUID2,
				Name:          "Liquid Oxygen",
				Description:   "Fuel",
				Price:         150.00,
				StockQuantity: 1000,
				Category:      entity.CATEGORY_FUEL,
			},
		}

		expectedTotalPrice = 25000150.00

	)

	s.inventoryClient.On("ListParts", s.ctx, filter).Return(partsFromInventory, nil)

	s.orderRepository.On("Create", s.ctx, mock.MatchedBy(func(order *entity.Order) bool {
		return order.UserUUID == userUUID &&
			len(order.PartUUIDs) == 2 &&
			order.TotalPrice == expectedTotalPrice &&
			order.Status == vo.OrderStatusPENDINGPAYMENT
	})).Return(nil)

	order, err := s.service.Create(s.ctx, request)

	s.Require().NoError(err)
	s.Require().NotNil(order)
	s.Require().Equal(userUUID, order.UserUUID)
	s.Require().Equal(2, len(order.PartUUIDs))
	s.Require().Equal(expectedTotalPrice, order.TotalPrice)
	s.Require().Equal(vo.OrderStatusPENDINGPAYMENT, order.Status)
	s.Require().NotEmpty(order.OrderUUID)
}

func (s *ServiceSuite) TestCreateOrderInventoryServiceError() {
	var (
		userUUID  = gofakeit.UUID()
		partUUID1 = gofakeit.UUID()
		partUUID2 = gofakeit.UUID()

		request = &dto.CreateOrderRequest{
			UserUUID:  userUUID,
			PartUUIDs: []string{partUUID1, partUUID2},
		}

		filter = &entity.PartsFilter{
			Uuids: []string{partUUID1, partUUID2},
		}

		inventoryError = errors.New("inventory service unavailable")
	)

	s.inventoryClient.On("ListParts", s.ctx, filter).Return(nil, inventoryError)

	order, err := s.service.Create(s.ctx, request)

	s.Require().Error(err)
	s.Require().Nil(order)
	s.Require().Contains(err.Error(), "failed to get parts")
}
