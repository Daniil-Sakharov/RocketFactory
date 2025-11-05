package dto

import "github.com/Daniil-Sakharov/RocketFactory/order/internal/model/vo"

type CreateOrderRequest struct {
	UserUUID  string   // UUID пользователя
	PartUUIDs []string // Список UUID деталей
}

type PayOrderRequest struct {
	OrderUUID     string           // UUID заказа
	PaymentMethod vo.PaymentMethod // Метод оплаты
}

type GetOrderRequest struct {
	OrderUUID string // UUID заказа
}

type CancelOrderRequest struct {
	OrderUUID string // UUID заказа
}

type CreateOrderResponse struct {
	OrderUUID  string  // UUID созданного заказа
	TotalPrice float64 // Общая стоимость
}
