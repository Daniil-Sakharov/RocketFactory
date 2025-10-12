package model

import "errors"

var (
	// ErrEmptyOrderUUID - ошибка когда UUID заказа пустой
	ErrEmptyOrderUUID = errors.New("order UUID is empty")

	// ErrEmptyUserUUID - ошибка когда UUID пользователя пустой
	ErrEmptyUserUUID = errors.New("user UUID is empty")

	// ErrInvalidPaymentMethod - ошибка когда метод оплаты не указан
	ErrInvalidPaymentMethod = errors.New("invalid payment method")
)
