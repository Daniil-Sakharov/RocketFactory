package converter

import (
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/dto"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/vo"
	paymentv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/payment/v1"
)

// PaymentRequestFromProto конвертирует protobuf запрос в domain модель
func PaymentRequestFromProto(req *paymentv1.PayOrderRequest) *dto.PayOrderClientRequest {
	return &dto.PayOrderClientRequest{
		OrderUUID:     req.GetOrderUuid(),
		UserUUID:      req.GetUserUuid(),
		PaymentMethod: PaymentMethodFromProto(req.GetPaymentMethod()),
	}
}

// PaymentResponseToProto конвертирует domain ответ в protobuf
func PaymentResponseToProto(resp *dto.PayOrderClientResponse) *paymentv1.PayOrderResponse {
	return &paymentv1.PayOrderResponse{
		TransactionUuid: resp.TransactionUUID,
	}
}

func PaymentResponseFromProto(resp *paymentv1.PayOrderResponse) *dto.PayOrderClientResponse {
	return &dto.PayOrderClientResponse{
		TransactionUUID: resp.TransactionUuid,
	}
}

// PaymentMethodFromProto конвертирует protobuf enum в domain enum
func PaymentMethodFromProto(protoMethod paymentv1.PaymentMethod) vo.PaymentMethod {
	switch protoMethod {
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CARD:
		return vo.PaymentMethodCARD
	case paymentv1.PaymentMethod_PAYMENT_METHOD_SBP:
		return vo.PaymentMethodSBP
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return vo.PaymentMethodCREDITCARD
	case paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return vo.PaymentMethodINVESTORMONEY
	default:
		return vo.PaymentMethodUNKNOWN
	}
}

// PaymentMethodToProto конвертирует domain enum в protobuf enum
func PaymentMethodToProto(method vo.PaymentMethod) paymentv1.PaymentMethod {
	switch method {
	case vo.PaymentMethodCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case vo.PaymentMethodSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	case vo.PaymentMethodCREDITCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case vo.PaymentMethodINVESTORMONEY:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
