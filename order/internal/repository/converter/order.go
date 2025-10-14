// internal/repository/converter/order.go
package converter

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/domain"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/vo"
	repoModel "github.com/Daniil-Sakharov/RocketFactory/order/internal/repository/model"
)

func RepoOrderToDomainModel(order *repoModel.Order) *domain.Order {
	return &domain.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice.InexactFloat64(),
		TransactionUUID: order.TransactionUUID.String,
		PaymentMethod:   vo.PaymentMethod(order.PaymentMethod),
		Status:          vo.OrderStatus(order.Status),
	}
}

func DomainOrderToRepoModel(order *domain.Order) *repoModel.Order {
	// Конвертация TransactionUUID: "" → NULL
	var txUUID sql.NullString
	if order.TransactionUUID != "" {
		txUUID = sql.NullString{String: order.TransactionUUID, Valid: true}
	}

	return &repoModel.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUUIDs:       pq.StringArray(order.PartUUIDs),
		TotalPrice:      decimal.NewFromFloat(order.TotalPrice),
		TransactionUUID: txUUID,
		PaymentMethod:   string(order.PaymentMethod),
		Status:          string(order.Status),
	}
}
