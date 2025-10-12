package dto

import "github.com/Daniil-Sakharov/RocketFactory/order/internal/model/vo"

// ==================== REQUEST DTOs ====================

// CreateOrderRequest - DTO для создания заказа
type CreateOrderRequest struct {
	UserUUID  string   // UUID пользователя
	PartUUIDs []string // Список UUID деталей
}

// PayOrderRequest - DTO для оплаты заказа
type PayOrderRequest struct {
	OrderUUID     string           // UUID заказа
	PaymentMethod vo.PaymentMethod // Метод оплаты
}

// GetOrderRequest - DTO для получения заказа
type GetOrderRequest struct {
	OrderUUID string // UUID заказа
}

// CancelOrderRequest - DTO для отмены заказа
type CancelOrderRequest struct {
	OrderUUID string // UUID заказа
}

// ==================== RESPONSE DTOs ====================

// CreateOrderResponse - DTO ответа на создание заказа
type CreateOrderResponse struct {
	OrderUUID  string  // UUID созданного заказа
	TotalPrice float64 // Общая стоимость
}
