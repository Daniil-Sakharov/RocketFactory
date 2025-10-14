package order

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/domain"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/dto"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/vo"
)

func (s *ServiceSuite) TestGetOrderSuccess() {
	var (
		orderUUID       = gofakeit.UUID()
		userUUID        = gofakeit.UUID()
		partUUID1       = gofakeit.UUID()
		partUUID2       = gofakeit.UUID()
		transactionUUID = gofakeit.UUID()

		request = &dto.GetOrderRequest{
			OrderUUID: orderUUID,
		}

		expectedOrder = &domain.Order{
			OrderUUID:       orderUUID,
			UserUUID:        userUUID,
			PartUUIDs:       []string{partUUID1, partUUID2},
			TotalPrice:      10_000_000.00,
			TransactionUUID: transactionUUID,
			PaymentMethod:   vo.PaymentMethodCARD,
			Status:          vo.OrderStatusPAID,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(expectedOrder, nil)

	order, err := s.service.Get(s.ctx, request)

	s.Require().NoError(err)
	s.Require().NotNil(order)
	s.Require().Equal(expectedOrder, order)
	s.Require().Equal(orderUUID, order.OrderUUID)
	s.Require().Equal(userUUID, order.UserUUID)
	s.Require().Equal(vo.OrderStatusPAID, order.Status)
}

func (s *ServiceSuite) TestGetOrderNotFound() {
	var (
		orderUUID = gofakeit.UUID()

		request = &dto.GetOrderRequest{
			OrderUUID: orderUUID,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(nil, model.ErrOrderNotFound)

	order, err := s.service.Get(s.ctx, request)

	s.Require().Error(err)
	s.Require().Nil(order)
	s.Require().ErrorIs(err, model.ErrOrderNotFound)
}
