package domain

import (
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/vo"
)

// Order - доменная сущность заказа
type Order struct {
	OrderUUID       string           // UUID заказа
	UserUUID        string           // UUID пользователя
	PartUUIDs       []string         // Список UUID деталей
	TotalPrice      float64          // Общая стоимость заказа
	TransactionUUID string           // UUID транзакции (если оплачен)
	PaymentMethod   vo.PaymentMethod // Способ оплаты
	Status          vo.OrderStatus   // Статус заказа
}
