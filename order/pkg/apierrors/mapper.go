// order/pkg/apierrors/mapper.go
package apierrors

import (
	"errors"
	"log"
	"strings"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model"
	orderV1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/openapi/order/v1"
)

// MapToCreateOrderError маппит ошибки для CreateOrder
func MapToCreateOrderError(err error) orderV1.CreateOrderRes {
	if err == nil {
		return nil
	}

	// Валидация → 400
	if errors.Is(err, model.ErrEmptyUserUUID) ||
		errors.Is(err, model.ErrEmptyPartUUIDs) ||
		errors.Is(err, model.ErrInvalidPaymentMethod) {
		return &orderV1.ValidationError{
			Error:   "VALIDATION_ERROR",
			Message: err.Error(),
		}
	}

	// Not Found → 404
	if errors.Is(err, model.ErrOrderNotFound) ||
		errors.Is(err, model.ErrPartsNotFound) {
		return &orderV1.NotFoundError{
			Error:   "NOT_FOUND",
			Message: err.Error(),
		}
	}

	// Conflict → 409
	if errors.Is(err, model.ErrOrderAlreadyExist) ||
		errors.Is(err, model.ErrOrderAlreadyPaid) {
		return &orderV1.ConflictError{
			Error:   "CONFLICT",
			Message: err.Error(),
		}
	}

	// External service → 502
	if isExternalServiceError(err) {
		return &orderV1.BadGatewayError{
			Error:   "EXTERNAL_SERVICE_ERROR",
			Message: "Failed to communicate with external service",
		}
	}

	// Internal → 500
	log.Printf("Unhandled error in CreateOrder: %v", err)
	return &orderV1.InternalServerError{
		Error:   "INTERNAL_ERROR",
		Message: "An internal error occurred",
	}
}

// MapToGetOrderError маппит ошибки для GetOrder
func MapToGetOrderError(err error) orderV1.GetOrderRes {
	if err == nil {
		return nil
	}

	// Not Found → 404
	if errors.Is(err, model.ErrOrderNotFound) {
		return &orderV1.NotFoundError{
			Error:   "NOT_FOUND",
			Message: err.Error(),
		}
	}

	// Internal → 500
	log.Printf("Unhandled error in GetOrder: %v", err)
	return &orderV1.InternalServerError{
		Error:   "INTERNAL_ERROR",
		Message: "An internal error occurred",
	}
}

// MapToPayOrderError маппит ошибки для PayOrder
func MapToPayOrderError(err error) orderV1.PayOrderRes {
	if err == nil {
		return nil
	}

	// Валидация → 400
	if errors.Is(err, model.ErrInvalidPaymentMethod) {
		return &orderV1.ValidationError{
			Error:   "VALIDATION_ERROR",
			Message: err.Error(),
		}
	}

	// Not Found → 404
	if errors.Is(err, model.ErrOrderNotFound) {
		return &orderV1.NotFoundError{
			Error:   "NOT_FOUND",
			Message: err.Error(),
		}
	}

	// Conflict → 409
	if errors.Is(err, model.ErrOrderAlreadyPaid) {
		return &orderV1.ConflictError{
			Error:   "CONFLICT",
			Message: err.Error(),
		}
	}

	// External service → 502
	if isExternalServiceError(err) {
		return &orderV1.BadGatewayError{
			Error:   "EXTERNAL_SERVICE_ERROR",
			Message: "Failed to communicate with external service",
		}
	}

	// Internal → 500
	log.Printf("Unhandled error in PayOrder: %v", err)
	return &orderV1.InternalServerError{
		Error:   "INTERNAL_ERROR",
		Message: "An internal error occurred",
	}
}

// MapToCancelOrderError маппит ошибки для CancelOrder
func MapToCancelOrderError(err error) orderV1.CancelOrderRes {
	if err == nil {
		return nil
	}

	// Not Found → 404
	if errors.Is(err, model.ErrOrderNotFound) {
		return &orderV1.NotFoundError{
			Error:   "NOT_FOUND",
			Message: err.Error(),
		}
	}

	// Conflict → 409
	if errors.Is(err, model.ErrOrderAlreadyPaid) {
		return &orderV1.ConflictError{
			Error:   "CONFLICT",
			Message: err.Error(),
		}
	}

	// Internal → 500
	log.Printf("Unhandled error in CancelOrder: %v", err)
	return &orderV1.InternalServerError{
		Error:   "INTERNAL_ERROR",
		Message: "An internal error occurred",
	}
}

// isExternalServiceError проверяет, связана ли ошибка с внешним сервисом
func isExternalServiceError(err error) bool {
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "payment service") ||
		strings.Contains(errMsg, "inventory service") ||
		strings.Contains(errMsg, "failed to get parts") ||
		strings.Contains(errMsg, "connection refused")
}
