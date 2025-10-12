package dto

import "github.com/Daniil-Sakharov/RocketFactory/order/internal/model/vo"

// PayOrderClientRequest - запрос на оплату заказа (для Payment Client)
type PayOrderClientRequest struct {
	OrderUUID     string           // UUID заказа
	UserUUID      string           // UUID пользователя, который производит оплату
	PaymentMethod vo.PaymentMethod // Метод оплаты
}

// PayOrderClientResponse - ответ на оплату заказа (от Payment Client)
type PayOrderClientResponse struct {
	TransactionUUID string // UUID транзакции
}
