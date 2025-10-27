package order

import (
	"errors"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/domain"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/dto"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/vo"
)

func (s *ServiceSuite) TestPayOrderSuccess() {
	var (
		orderUUID       = gofakeit.UUID()
		userUUID        = gofakeit.UUID()
		partUUID1       = gofakeit.UUID()
		partUUID2       = gofakeit.UUID()
		partsUUIDs      = []string{partUUID1, partUUID2}
		expectedPrice   = 20_000.00
		paymentMethod   = vo.PaymentMethodCARD
		transactionUUID = gofakeit.UUID()

		payOrderRequest = &dto.PayOrderRequest{
			OrderUUID:     orderUUID,
			PaymentMethod: paymentMethod,
		}

		payOrderClientRequest = &dto.PayOrderClientRequest{
			OrderUUID:     orderUUID,
			UserUUID:      userUUID,
			PaymentMethod: paymentMethod,
		}

		payOrderClientResponse = &dto.PayOrderClientResponse{
			TransactionUUID: transactionUUID,
		}

		orderFromDB = &domain.Order{
			OrderUUID:       orderUUID,
			UserUUID:        userUUID,
			PartUUIDs:       partsUUIDs,
			TotalPrice:      expectedPrice,
			TransactionUUID: "",
			PaymentMethod:   "",
			Status:          vo.OrderStatusPENDINGPAYMENT,
		}

		expectedUpdatedOrder = &domain.Order{
			OrderUUID:       orderUUID,
			UserUUID:        userUUID,
			PartUUIDs:       partsUUIDs,
			TotalPrice:      expectedPrice,
			TransactionUUID: transactionUUID,
			PaymentMethod:   paymentMethod,
			Status:          vo.OrderStatusPAID,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(orderFromDB, nil)

	s.paymentClient.On("PayOrder", s.ctx, payOrderClientRequest).Return(payOrderClientResponse, nil)

	s.orderRepository.On("Update", s.ctx, mock.MatchedBy(func(order *domain.Order) bool {
		return order.OrderUUID == orderUUID &&
			order.TransactionUUID == transactionUUID &&
			order.PaymentMethod == paymentMethod &&
			order.Status == vo.OrderStatusPAID
	})).Return(nil)

	s.orderProducer.On("PublishOrderPaid", s.ctx, mock.Anything).Return(nil)

	order, err := s.service.Pay(s.ctx, payOrderRequest)

	s.Require().NoError(err)
	s.Require().NotNil(order)
	s.Require().Equal(expectedUpdatedOrder.OrderUUID, order.OrderUUID)
	s.Require().Equal(expectedUpdatedOrder.TransactionUUID, order.TransactionUUID)
	s.Require().Equal(expectedUpdatedOrder.PaymentMethod, order.PaymentMethod)
	s.Require().Equal(vo.OrderStatusPAID, order.Status)
}

func (s *ServiceSuite) TestPayOrderNotFound() {
	var (
		orderUUID = gofakeit.UUID()

		payOrderRequest = &dto.PayOrderRequest{
			OrderUUID:     orderUUID,
			PaymentMethod: vo.PaymentMethodCARD,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(nil, model.ErrOrderNotFound)

	order, err := s.service.Pay(s.ctx, payOrderRequest)

	s.Require().Error(err)
	s.Require().Nil(order)
	s.Require().ErrorIs(err, model.ErrOrderNotFound)
}

func (s *ServiceSuite) TestPayOrderPaymentClientError() {
	var (
		orderUUID     = gofakeit.UUID()
		userUUID      = gofakeit.UUID()
		partUUID1     = gofakeit.UUID()
		partUUID2     = gofakeit.UUID()
		partsUUIDs    = []string{partUUID1, partUUID2}
		expectedPrice = 20_000.00
		paymentMethod = vo.PaymentMethodCARD

		payOrderRequest = &dto.PayOrderRequest{
			OrderUUID:     orderUUID,
			PaymentMethod: paymentMethod,
		}

		payOrderClientRequest = &dto.PayOrderClientRequest{
			OrderUUID:     orderUUID,
			UserUUID:      userUUID,
			PaymentMethod: paymentMethod,
		}

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

	s.paymentClient.On("PayOrder", s.ctx, payOrderClientRequest).Return(nil, errors.New("payment service unavailable"))

	order, err := s.service.Pay(s.ctx, payOrderRequest)

	s.Require().Error(err)
	s.Require().Nil(order)
	s.Require().Contains(err.Error(), "failed to access pay")
}
