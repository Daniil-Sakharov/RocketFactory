package app

import (
	"context"
	api "github.com/Daniil-Sakharov/RocketFactory/payment/internal/api/payment/v1"
	"github.com/Daniil-Sakharov/RocketFactory/payment/internal/service"
	"github.com/Daniil-Sakharov/RocketFactory/payment/internal/service/payment"
	paymentv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentV1API   paymentv1.PaymentServiceServer
	paymentService service.PaymentService
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentAPI(ctx context.Context) paymentv1.PaymentServiceServer {
	if d.paymentV1API == nil {
		d.paymentV1API = api.New(d.PaymentService(ctx))
	}
	return d.paymentV1API
}

func (d *diContainer) PaymentService(_ context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = payment.New()
	}
	return d.paymentService
}
