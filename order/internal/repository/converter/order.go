package converter

import (
	"github.com/google/uuid"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/entity"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/vo"
	repoModel "github.com/Daniil-Sakharov/RocketFactory/order/internal/repository/model"
	"github.com/Daniil-Sakharov/RocketFactory/order/pkg/utils"
)

// OrderToRepoModel конвертирует Domain Order → Repository Order
func OrderToRepoModel(order *entity.Order) *repoModel.Order {
	return &repoModel.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        uuid.MustParse(order.UserUUID),
		PartUUIDs:       utils.StringsToUUIDs(order.PartUUIDs),
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   PaymentMethodToRepoModel(order.PaymentMethod),
		Status:          StatusToRepoModel(order.Status),
	}
}

// OrderFromRepoModel конвертирует Repository Order → Domain Order
func OrderFromRepoModel(repoOrder *repoModel.Order) *entity.Order {
	return &entity.Order{
		OrderUUID:       repoOrder.OrderUUID,
		UserUUID:        repoOrder.UserUUID.String(),
		PartUUIDs:       utils.UuidsToStrings(repoOrder.PartUUIDs),
		TotalPrice:      repoOrder.TotalPrice,
		TransactionUUID: repoOrder.TransactionUUID,
		PaymentMethod:   PaymentMethodFromRepoModel(repoOrder.PaymentMethod),
		Status:          StatusFromRepoModel(repoOrder.Status),
	}
}

// StatusToRepoModel конвертирует Domain OrderStatus → Repository OrderStatus
func StatusToRepoModel(status vo.OrderStatus) repoModel.OrderStatus {
	switch status {
	case vo.OrderStatusPAID:
		return repoModel.OrderStatusPAID
	case vo.OrderStatusPENDINGPAYMENT:
		return repoModel.OrderStatusPENDINGPAYMENT
	case vo.OrderStatusCANCELLED:
		return repoModel.OrderStatusCANCELLED
	default:
		return repoModel.OrderStatusCANCELLED
	}
}

// StatusFromRepoModel конвертирует Repository OrderStatus → Domain OrderStatus
func StatusFromRepoModel(status repoModel.OrderStatus) vo.OrderStatus {
	switch status {
	case repoModel.OrderStatusPAID:
		return vo.OrderStatusPAID
	case repoModel.OrderStatusPENDINGPAYMENT:
		return vo.OrderStatusPENDINGPAYMENT
	case repoModel.OrderStatusCANCELLED:
		return vo.OrderStatusCANCELLED
	default:
		return vo.OrderStatusCANCELLED
	}
}

// PaymentMethodToRepoModel конвертирует Domain PaymentMethod → Repository PaymentMethod
func PaymentMethodToRepoModel(paymentMethod vo.PaymentMethod) repoModel.PaymentMethod {
	switch paymentMethod {
	case vo.PaymentMethodCARD:
		return repoModel.PaymentMethodCARD
	case vo.PaymentMethodCREDITCARD:
		return repoModel.PaymentMethodCREDITCARD
	case vo.PaymentMethodINVESTORMONEY:
		return repoModel.PaymentMethodINVESTORMONEY
	case vo.PaymentMethodSBP:
		return repoModel.PaymentMethodSBP
	case vo.PaymentMethodUNKNOWN:
		return repoModel.PaymentMethodUNKNOWN
	default:
		return repoModel.PaymentMethodUNKNOWN
	}
}

// PaymentMethodFromRepoModel конвертирует Repository PaymentMethod → Domain PaymentMethod
func PaymentMethodFromRepoModel(paymentMethod repoModel.PaymentMethod) vo.PaymentMethod {
	switch paymentMethod {
	case repoModel.PaymentMethodCARD:
		return vo.PaymentMethodCARD
	case repoModel.PaymentMethodCREDITCARD:
		return vo.PaymentMethodCREDITCARD
	case repoModel.PaymentMethodINVESTORMONEY:
		return vo.PaymentMethodINVESTORMONEY
	case repoModel.PaymentMethodSBP:
		return vo.PaymentMethodSBP
	case repoModel.PaymentMethodUNKNOWN:
		return vo.PaymentMethodUNKNOWN
	default:
		return vo.PaymentMethodUNKNOWN
	}
}
