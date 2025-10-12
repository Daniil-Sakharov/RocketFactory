package payment

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/Daniil-Sakharov/RocketFactory/payment/internal/model"
)

func (s *ServiceSuite) TestPayOrderSuccess() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()

		request = &model.PayOrderRequest{
			OrderUUID:     orderUUID,
			UserUUID:      userUUID,
			PaymentMethod: model.PaymentMethodCard,
		}
	)

	response, err := s.service.PayOrder(s.ctx, request)

	s.Require().NoError(err)
	s.Require().NotNil(response)
	s.Require().NotEmpty(response.TransactionUUID)
}

func (s *ServiceSuite) TestPayOrderEmptyOrderUUID() {
	var (
		userUUID = gofakeit.UUID()

		request = &model.PayOrderRequest{
			OrderUUID:     "",
			UserUUID:      userUUID,
			PaymentMethod: model.PaymentMethodCard,
		}
	)

	response, err := s.service.PayOrder(s.ctx, request)

	s.Require().Error(err)
	s.Require().Nil(response)
	s.Require().ErrorIs(err, model.ErrEmptyOrderUUID)
}

func (s *ServiceSuite) TestPayOrderEmptyUserUUID() {
	var (
		orderUUID = gofakeit.UUID()

		request = &model.PayOrderRequest{
			OrderUUID:     orderUUID,
			UserUUID:      "",
			PaymentMethod: model.PaymentMethodCard,
		}
	)

	response, err := s.service.PayOrder(s.ctx, request)

	s.Require().Error(err)
	s.Require().Nil(response)
	s.Require().ErrorIs(err, model.ErrEmptyUserUUID)
}

func (s *ServiceSuite) TestPayOrderInvalidPaymentMethod() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()

		request = &model.PayOrderRequest{
			OrderUUID:     orderUUID,
			UserUUID:      userUUID,
			PaymentMethod: model.PaymentMethodUnspecified,
		}
	)

	response, err := s.service.PayOrder(s.ctx, request)

	s.Require().Error(err)
	s.Require().Nil(response)
	s.Require().ErrorIs(err, model.ErrInvalidPaymentMethod)
}

func (s *ServiceSuite) TestPayOrderWithSBP() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()

		request = &model.PayOrderRequest{
			OrderUUID:     orderUUID,
			UserUUID:      userUUID,
			PaymentMethod: model.PaymentMethodSBP,
		}
	)

	response, err := s.service.PayOrder(s.ctx, request)

	s.Require().NoError(err)
	s.Require().NotNil(response)
	s.Require().NotEmpty(response.TransactionUUID)
}

func (s *ServiceSuite) TestPayOrderWithCreditCard() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()

		request = &model.PayOrderRequest{
			OrderUUID:     orderUUID,
			UserUUID:      userUUID,
			PaymentMethod: model.PaymentMethodCreditCard,
		}
	)

	response, err := s.service.PayOrder(s.ctx, request)

	s.Require().NoError(err)
	s.Require().NotNil(response)
	s.Require().NotEmpty(response.TransactionUUID)
}

func (s *ServiceSuite) TestPayOrderWithInvestorMoney() {
	var (
		orderUUID = gofakeit.UUID()
		userUUID  = gofakeit.UUID()

		request = &model.PayOrderRequest{
			OrderUUID:     orderUUID,
			UserUUID:      userUUID,
			PaymentMethod: model.PaymentMethodInvestorMoney,
		}
	)

	response, err := s.service.PayOrder(s.ctx, request)

	s.Require().NoError(err)
	s.Require().NotNil(response)
	s.Require().NotEmpty(response.TransactionUUID)
}
