package v1

import (
	"github.com/Daniil-Sakharov/RocketFactory/payment/internal/service"
	paymentv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/payment/v1"
)

// api - реализация gRPC сервера для PaymentService
type api struct {
	paymentv1.UnimplementedPaymentServiceServer
	paymentService service.PaymentService
}

// New создает новый экземпляр API
func New(paymentService service.PaymentService) *api {
	return &api{
		paymentService: paymentService,
	}
}
