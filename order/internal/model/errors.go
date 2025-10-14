package model

import "errors"

var (
	ErrOrderAlreadyExist     = errors.New("order already exist")
	ErrOrderNotFound         = errors.New("order not found")
	ErrOrderAlreadyPaid      = errors.New("order already paid")
	ErrOrderAlreadyCancelled = errors.New("order already cancelled")
	ErrEmptyUserUUID         = errors.New("user UUID is empty")
	ErrEmptyPartUUIDs        = errors.New("part UUIDs are empty")
	ErrPartsNotFound         = errors.New("parts not found")
	ErrInvalidPaymentMethod  = errors.New("invalid payment method")
	ErrUnknownError          = errors.New("unknown error")
)
