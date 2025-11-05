package dto

import "github.com/Daniil-Sakharov/RocketFactory/order/internal/model/vo"

type PayOrderClientRequest struct {
	OrderUUID     string           // UUID заказа
	UserUUID      string           // UUID пользователя, который производит оплату
	PaymentMethod vo.PaymentMethod // Метод оплаты
}

type PayOrderClientResponse struct {
	TransactionUUID string // UUID транзакции
}
