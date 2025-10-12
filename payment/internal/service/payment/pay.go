package payment

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/Daniil-Sakharov/RocketFactory/payment/internal/model"
)

// PayOrder обрабатывает платеж заказа
func (s *svc) PayOrder(ctx context.Context, req *model.PayOrderRequest) (*model.PayOrderResponse, error) {
	// 1. Валидация входных данных
	if err := s.validatePaymentRequest(req); err != nil {
		return nil, err
	}

	// 2. Генерация UUID транзакции
	transactionUUID := uuid.NewString()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUUID)

	// 3. Здесь могла бы быть логика:
	// - Вызов платежного шлюза (Stripe, PayPal)
	// - Проверка баланса пользователя
	// - Сохранение транзакции в БД
	// - Отправка события "PaymentProcessed"

	// 4. Возврат результата
	return &model.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}

// validatePaymentRequest проверяет корректность запроса на оплату
func (s *svc) validatePaymentRequest(req *model.PayOrderRequest) error {
	if req.OrderUUID == "" {
		return model.ErrEmptyOrderUUID
	}

	if req.UserUUID == "" {
		return model.ErrEmptyUserUUID
	}

	if req.PaymentMethod == model.PaymentMethodUnspecified {
		return model.ErrInvalidPaymentMethod
	}

	return nil
}
