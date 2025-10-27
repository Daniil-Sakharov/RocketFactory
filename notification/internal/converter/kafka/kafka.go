package kafka

import "github.com/Daniil-Sakharov/RocketFactory/notification/internal/model/domain"

type OrderDecoder interface {
	OrderDecode(data []byte) (domain.OrderConsumeEvent, error)
}

type AssemblyDecoder interface {
	AssemblyDecode(data []byte) (domain.AssemblyConsumeEvent, error)
}
