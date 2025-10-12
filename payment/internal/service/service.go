package service

import (
	"context"

	"github.com/Daniil-Sakharov/RocketFactory/payment/internal/model"
)

type PaymentService interface {
	// PayOrder обрабатывает платеж и возвращает UUID транзакции
	PayOrder(ctx context.Context, req *model.PayOrderRequest) (*model.PayOrderResponse, error)
}
