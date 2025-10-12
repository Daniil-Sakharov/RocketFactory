package vo

// PaymentMethod - способ оплаты
type PaymentMethod string

const (
	// PaymentMethodUNKNOWN - неизвестный способ оплаты
	PaymentMethodUNKNOWN PaymentMethod = "UNKNOWN"
	// PaymentMethodCARD - банковская карта
	PaymentMethodCARD PaymentMethod = "CARD"
	// PaymentMethodSBP - Система Быстрых Платежей
	PaymentMethodSBP PaymentMethod = "SBP"
	// PaymentMethodCREDITCARD - кредитная карта
	PaymentMethodCREDITCARD PaymentMethod = "CREDIT_CARD"
	// PaymentMethodINVESTORMONEY - деньги инвестора (внутренний метод)
	PaymentMethodINVESTORMONEY PaymentMethod = "INVESTOR_MONEY"
)
