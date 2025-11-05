package order

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/domain"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/vo"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/service/dto"
)

func (s *ServiceSuite) TestCancelOrderSuccess() {
	var (
		orderUUID     = gofakeit.UUID()
		userUUID      = gofakeit.UUID()
		partUUID1     = gofakeit.UUID()
		partUUID2     = gofakeit.UUID()
		partsUUIDs    = []string{partUUID1, partUUID2}
		expectedPrice = 20_000.00

		cancelOrderRequest = &dto.CancelOrderRequest{OrderUUID: orderUUID}

		orderFromDB = &domain.Order{
			OrderUUID:       orderUUID,
			UserUUID:        userUUID,
			PartUUIDs:       partsUUIDs,
			TotalPrice:      expectedPrice,
			TransactionUUID: "",
			PaymentMethod:   "",
			Status:          vo.OrderStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(orderFromDB, nil)

	s.orderRepository.On("Update", s.ctx, mock.MatchedBy(func(order *domain.Order) bool {
		return order.OrderUUID == orderUUID &&
			order.UserUUID == userUUID &&
			order.TransactionUUID == "" &&
			order.PaymentMethod == "" &&
			order.Status == vo.OrderStatusCANCELLED
	})).Return(nil)

	err := s.service.Cancel(s.ctx, cancelOrderRequest)

	s.Require().NoError(err)
}

func (s *ServiceSuite) TestCancelOrderAlreadyPaid() {
	var (
		orderUUID       = gofakeit.UUID()
		userUUID        = gofakeit.UUID()
		partUUID1       = gofakeit.UUID()
		partUUID2       = gofakeit.UUID()
		partsUUIDs      = []string{partUUID1, partUUID2}
		transactionUUID = gofakeit.UUID()

		cancelOrderRequest = &dto.CancelOrderRequest{
			OrderUUID: orderUUID,
		}

		paidOrderFromDB = &domain.Order{
			OrderUUID:       orderUUID,
			UserUUID:        userUUID,
			PartUUIDs:       partsUUIDs,
			TotalPrice:      20_000.00,
			TransactionUUID: transactionUUID,
			PaymentMethod:   vo.PaymentMethodCARD,
			Status:          vo.OrderStatusPAID,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(paidOrderFromDB, nil)

	err := s.service.Cancel(s.ctx, cancelOrderRequest)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderAlreadyPaid)
}

func (s *ServiceSuite) TestCancelOrderNotFound() {
	var (
		orderUUID = gofakeit.UUID()

		cancelOrderRequest = &dto.CancelOrderRequest{
			OrderUUID: orderUUID,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(nil, model.ErrOrderNotFound)

	err := s.service.Cancel(s.ctx, cancelOrderRequest)

	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrOrderNotFound)
}
