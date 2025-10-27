package kafka

import "github.com/Daniil-Sakharov/RocketFactory/assembly/internal/model"

// OrderPaidDecoder - декодер для OrderPaid события
type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaidEvent, error)
}
