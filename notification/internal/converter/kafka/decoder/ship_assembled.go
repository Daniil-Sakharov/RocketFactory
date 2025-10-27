package decoder

import (
	"fmt"
	def "github.com/Daniil-Sakharov/RocketFactory/notification/internal/converter/kafka"
	"github.com/Daniil-Sakharov/RocketFactory/notification/internal/model/domain"
	eventsv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/events/v1"
	"google.golang.org/protobuf/proto"
	"time"
)

var _ def.AssemblyDecoder = (*assemblyDecoder)(nil)

type assemblyDecoder struct{}

func NewAssemblyDecoder() *assemblyDecoder {
	return &assemblyDecoder{}
}

func (d *assemblyDecoder) AssemblyDecode(data []byte) (domain.AssemblyConsumeEvent, error) {
	var pb eventsv1.ShipAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return domain.AssemblyConsumeEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return domain.AssemblyConsumeEvent{
		EventUUID:       pb.EventUuid,
		OrderUUID:       pb.OrderUuid,
		UserUUID:        pb.UserUuid,
		BuildTimeSec: time.Duration(pb.BuildTimeSec),
	}, nil
}

