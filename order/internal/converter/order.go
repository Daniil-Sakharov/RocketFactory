package converter

import (
	"github.com/google/uuid"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/dto"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/entity"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/vo"
	"github.com/Daniil-Sakharov/RocketFactory/order/pkg/utils"
	orderV1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/openapi/order/v1"
)

func CreateOrderRequestToServiceModel(req orderV1.CreateOrderRequest) *dto.CreateOrderRequest {
	return &dto.CreateOrderRequest{
		UserUUID:  req.UserUUID.String(),
		PartUUIDs: utils.UuidsToStrings(req.PartUuids),
	}
}

func CreateOrderResponseFromEntity(order *entity.Order) *orderV1.CreateOrderResponse {
	orderUUID := uuid.MustParse(order.OrderUUID)

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: order.TotalPrice,
	}
}

func GetOrderResponseFromEntity(order *entity.Order) *orderV1.GetOrderResponse {
	orderUUID := uuid.MustParse(order.OrderUUID)
	userUUID := uuid.MustParse(order.UserUUID)

	partUUIDs := make([]uuid.UUID, 0, len(order.PartUUIDs))
	for _, partUUIDStr := range order.PartUUIDs {
		if partUUID, err := uuid.Parse(partUUIDStr); err == nil {
			partUUIDs = append(partUUIDs, partUUID)
		}
	}

	var transactionUUID orderV1.OptUUID
	if order.TransactionUUID != "" {
		if txUUID, err := uuid.Parse(order.TransactionUUID); err == nil {
			transactionUUID.SetTo(txUUID)
		}
	}

	var paymentMethod orderV1.OptPaymentMethod
	if order.PaymentMethod != vo.PaymentMethodUNKNOWN {
		paymentMethod.SetTo(PaymentMethodToOpenAPI(order.PaymentMethod))
	}

	return &orderV1.GetOrderResponse{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      order.TotalPrice,
		Status:          OrderStatusToOpenAPI(order.Status),
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
	}
}

func PayOrderRequestToServiceModel(req orderV1.PayOrderRequest, orderUUID string) *dto.PayOrderRequest {
	return &dto.PayOrderRequest{
		OrderUUID:     orderUUID,
		PaymentMethod: PaymentMethodFromOpenAPI(req.PaymentMethod),
	}
}

func GetOrderRequestToServiceModel(orderUUID string) *dto.GetOrderRequest {
	return &dto.GetOrderRequest{
		OrderUUID: orderUUID,
	}
}

func CancelOrderRequestToServiceModel(orderUUID string) *dto.CancelOrderRequest {
	return &dto.CancelOrderRequest{
		OrderUUID: orderUUID,
	}
}

func PayOrderResponseFromEntity(order *entity.Order) *orderV1.PayOrderResponse {
	transactionUUID := uuid.MustParse(order.TransactionUUID)

	return &orderV1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}
}

func OrderStatusToOpenAPI(status vo.OrderStatus) orderV1.OrderStatus {
	switch status {
	case vo.OrderStatusPENDINGPAYMENT:
		return orderV1.OrderStatusPENDINGPAYMENT
	case vo.OrderStatusPAID:
		return orderV1.OrderStatusPAID
	case vo.OrderStatusCANCELLED:
		return orderV1.OrderStatusCANCELLED
	default:
		return orderV1.OrderStatusPENDINGPAYMENT
	}
}

func OrderStatusFromOpenAPI(status orderV1.OrderStatus) vo.OrderStatus {
	switch status {
	case orderV1.OrderStatusPENDINGPAYMENT:
		return vo.OrderStatusPENDINGPAYMENT
	case orderV1.OrderStatusPAID:
		return vo.OrderStatusPAID
	case orderV1.OrderStatusCANCELLED:
		return vo.OrderStatusCANCELLED
	default:
		return vo.OrderStatusPENDINGPAYMENT
	}
}

func PaymentMethodToOpenAPI(method vo.PaymentMethod) orderV1.PaymentMethod {
	switch method {
	case vo.PaymentMethodCARD:
		return orderV1.PaymentMethodCARD
	case vo.PaymentMethodSBP:
		return orderV1.PaymentMethodSBP
	case vo.PaymentMethodCREDITCARD:
		return orderV1.PaymentMethodCREDITCARD
	case vo.PaymentMethodINVESTORMONEY:
		return orderV1.PaymentMethodINVESTORMONEY
	default:
		return orderV1.PaymentMethodUNKNOWN
	}
}

func PaymentMethodFromOpenAPI(method orderV1.PaymentMethod) vo.PaymentMethod {
	switch method {
	case orderV1.PaymentMethodCARD:
		return vo.PaymentMethodCARD
	case orderV1.PaymentMethodSBP:
		return vo.PaymentMethodSBP
	case orderV1.PaymentMethodCREDITCARD:
		return vo.PaymentMethodCREDITCARD
	case orderV1.PaymentMethodINVESTORMONEY:
		return vo.PaymentMethodINVESTORMONEY
	default:
		return vo.PaymentMethodUNKNOWN
	}
}
