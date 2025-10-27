package domain

import "time"

type OrderProduceEvent struct {
	EventUUID       string
	OrderUUID       string
	UserUUID        string
	PaymentMethod   string
	TransactionUUID string
}

type AssemblyConsumeEvent struct {
	EventUUID string
	OrderUUID string
	UserUUID  string
	BuildTime time.Duration
}
