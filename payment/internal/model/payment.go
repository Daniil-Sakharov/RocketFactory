package model

// PaymentMethod - способ оплаты
type PaymentMethod int32

const (
	PaymentMethodUnspecified   PaymentMethod = 0 // Неизвестный способ
	PaymentMethodCard          PaymentMethod = 1 // Банковская карта
	PaymentMethodSBP           PaymentMethod = 2 // Система Быстрых Платежей
	PaymentMethodCreditCard    PaymentMethod = 3 // Кредитная карта
	PaymentMethodInvestorMoney PaymentMethod = 4 // Деньги инвестора (внутренний метод)
)

// PayOrderRequest - запрос на оплату заказа
type PayOrderRequest struct {
	OrderUUID     string        // UUID заказа
	UserUUID      string        // UUID пользователя, который производит оплату
	PaymentMethod PaymentMethod // Метод оплаты
}

// PayOrderResponse - ответ на оплату заказа
type PayOrderResponse struct {
	TransactionUUID string // UUID транзакции
}
