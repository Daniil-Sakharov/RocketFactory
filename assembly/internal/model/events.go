package model

import "time"

// ShipAssembledEvent - событие "корабль собран"
type ShipAssembledEvent struct {
	EventUUID    string
	OrderUUID    string
	UserUUID     string
	BuildTimeSec time.Duration
}

// OrderPaidEvent - событие "заказ оплачен" (приходит от Order Service)
type OrderPaidEvent struct {
	EventUUID       string
	OrderUUID       string
	UserUUID        string
	PaymentMethod   string
	TransactionUUID string
}
