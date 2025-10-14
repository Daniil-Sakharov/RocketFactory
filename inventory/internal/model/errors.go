package model

import "errors"

// ErrPartNotFound возвращается когда деталь не найдена
var (
	ErrPartNotFound  = errors.New("part not found")
	ErrPartsNotFound = errors.New("parts not found")
)
