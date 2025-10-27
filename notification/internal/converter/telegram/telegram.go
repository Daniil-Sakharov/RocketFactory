package telegram

import (
	"github.com/Daniil-Sakharov/RocketFactory/notification/internal/model/domain"
)

func OrderPaidEventToTemplateData(event *domain.OrderConsumeEvent) *domain.OrderTemplateData {
	return &domain.OrderTemplateData{
		OrderUUID:       event.OrderUUID,
		UserUUID:        event.UserUUID,
		PaymentMethod:   event.PaymentMethod,
		TransactionUUID: event.TransactionUUID,
	}
}

func ShipAssembledEventToTemplateData(event *domain.AssemblyConsumeEvent) *domain.AssembledTemplateData {
	return &domain.AssembledTemplateData{
		OrderUUID:    event.OrderUUID,
		UserUUID:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	}
}
