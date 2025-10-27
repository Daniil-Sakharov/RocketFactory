package domain

import "time"

type AssembledTemplateData struct{
	OrderUUID string
	UserUUID string
	BuildTimeSec time.Duration
}

type OrderTemplateData struct{
	OrderUUID string
	UserUUID string
	PaymentMethod string
	TransactionUUID string
}
