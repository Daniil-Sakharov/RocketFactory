package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Daniil-Sakharov/RocketFactory/payment/internal/converter"
	"github.com/Daniil-Sakharov/RocketFactory/payment/internal/model"
	paymentv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/payment/v1"
)

// PayOrder обрабатывает gRPC запрос на оплату заказа
func (a *api) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	// 1. Конвертация protobuf → domain
	paymentReq := converter.PaymentRequestFromProto(req)

	// 2. Вызов сервисного слоя
	paymentResp, err := a.paymentService.PayOrder(ctx, paymentReq)
	if err != nil {
		// Обработка domain ошибок и конвертация в gRPC статусы
		if errors.Is(err, model.ErrEmptyOrderUUID) ||
			errors.Is(err, model.ErrEmptyUserUUID) ||
			errors.Is(err, model.ErrInvalidPaymentMethod) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	// 3. Конвертация domain → protobuf
	protoResp := converter.PaymentResponseToProto(paymentResp)

	return protoResp, nil
}
