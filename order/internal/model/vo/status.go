package vo

// OrderStatus - статус заказа
type OrderStatus string

const (
	// OrderStatusPENDINGPAYMENT - заказ ожидает оплаты
	OrderStatusPENDINGPAYMENT OrderStatus = "PENDING_PAYMENT"
	// OrderStatusPAID - заказ оплачен
	OrderStatusPAID OrderStatus = "PAID"
	// OrderStatusCANCELLED - заказ отменен
	OrderStatusCANCELLED OrderStatus = "CANCELLED"
)
