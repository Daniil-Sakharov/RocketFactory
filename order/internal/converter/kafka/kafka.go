package kafka

import "github.com/Daniil-Sakharov/RocketFactory/order/internal/model/domain"

type AssemblyDecoder interface {
	Decode(data []byte) (domain.AssemblyConsumeEvent, error)
}
