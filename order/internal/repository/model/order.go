package model

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type Order struct {
	OrderUUID       string          `db:"order_uuid"`
	UserUUID        string          `db:"user_uuid"`
	PartUUIDs       pq.StringArray  `db:"part_uuids"`
	TotalPrice      decimal.Decimal `db:"total_price"`
	TransactionUUID sql.NullString  `db:"transaction_uuid"`
	PaymentMethod   string          `db:"payment_method"`
	Status          string          `db:"order_status"`
	CreatedAt       time.Time       `db:"created_at"`
	UpdatedAt       time.Time       `db:"updated_at"`
}

type OrderStatus string

const (
	OrderStatusPENDINGPAYMENT OrderStatus = "PENDING_PAYMENT"
	OrderStatusPAID           OrderStatus = "PAID"
	OrderStatusCANCELLED      OrderStatus = "CANCELLED"
)

type PaymentMethod string

const (
	PaymentMethodUNKNOWN       PaymentMethod = "UNKNOWN"
	PaymentMethodCARD          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCREDITCARD    PaymentMethod = "CREDIT_CARD"
	PaymentMethodINVESTORMONEY PaymentMethod = "INVESTOR_MONEY"
)
