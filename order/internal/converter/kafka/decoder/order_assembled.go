package decoder

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	def "github.com/Daniil-Sakharov/RocketFactory/order/internal/converter/kafka"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/model/domain"
	eventsv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/events/v1"
)

var _ def.AssemblyDecoder = (*decoder)(nil)

type decoder struct{}

func NewAssemblyDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (domain.AssemblyConsumeEvent, error) {
	var pb eventsv1.ShipAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return domain.AssemblyConsumeEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return domain.AssemblyConsumeEvent{
		EventUUID: pb.EventUuid,
		OrderUUID: pb.OrderUuid,
		UserUUID:  pb.UserUuid,
		BuildTime: time.Duration(pb.BuildTimeSec) * time.Second,
	}, nil
}
