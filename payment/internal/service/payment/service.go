package payment

import (
	"github.com/Daniil-Sakharov/RocketFactory/payment/internal/service"
)

// Проверка, что service реализует интерфейс PaymentService на этапе компиляции
var _ service.PaymentService = (*svc)(nil)

// svc - реализация PaymentService
type svc struct {
	// Пока пустая структура
	// В будущем можно добавить:
	// - logger Logger
	// - metrics MetricsCollector
	// - paymentGateway PaymentGateway
}

// New создает новый экземпляр PaymentService
func New() *svc {
	return &svc{}
}
