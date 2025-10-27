package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	def "github.com/Daniil-Sakharov/RocketFactory/notification/internal/converter/kafka"
	"github.com/Daniil-Sakharov/RocketFactory/notification/internal/model/domain"
	eventsv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/events/v1"
)

var _ def.OrderDecoder = (*orderDecoder)(nil)

type orderDecoder struct{}

func NewOrderDecoder() *orderDecoder {
	return &orderDecoder{}
}

func (d *orderDecoder) OrderDecode(data []byte) (domain.OrderConsumeEvent, error) {
	var pb eventsv1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return domain.OrderConsumeEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return domain.OrderConsumeEvent{
		EventUUID:       pb.EventUuid,
		OrderUUID:       pb.OrderUuid,
		UserUUID:        pb.UserUuid,
		PaymentMethod:   pb.PaymentMethod,
		TransactionUUID: pb.TransactionUuid,
	}, nil
}
