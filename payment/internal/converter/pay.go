package converter

import (
	"github.com/Daniil-Sakharov/RocketFactory/payment/internal/model"
	paymentv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/payment/v1"
)

// PaymentRequestFromProto конвертирует protobuf запрос в domain модель
func PaymentRequestFromProto(req *paymentv1.PayOrderRequest) *model.PayOrderRequest {
	return &model.PayOrderRequest{
		OrderUUID:     req.GetOrderUuid(),
		UserUUID:      req.GetUserUuid(),
		PaymentMethod: PaymentMethodFromProto(req.GetPaymentMethod()),
	}
}

// PaymentResponseToProto конвертирует domain ответ в protobuf
func PaymentResponseToProto(resp *model.PayOrderResponse) *paymentv1.PayOrderResponse {
	return &paymentv1.PayOrderResponse{
		TransactionUuid: resp.TransactionUUID,
	}
}

// PaymentMethodFromProto конвертирует protobuf enum в domain enum
func PaymentMethodFromProto(protoMethod paymentv1.PaymentMethod) model.PaymentMethod {
	switch protoMethod {
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case paymentv1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnspecified
	}
}

// PaymentMethodToProto конвертирует domain enum в protobuf enum
func PaymentMethodToProto(method model.PaymentMethod) paymentv1.PaymentMethod {
	switch method {
	case model.PaymentMethodCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCreditCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodInvestorMoney:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
