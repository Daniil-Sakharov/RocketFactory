package service

import (
	"context"
	"github.com/Daniil-Sakharov/RocketFactory/notification/internal/model/domain"
)

type TelegramService interface {
	SendShipAssembledNotification(ctx context.Context, templateData *domain.AssembledTemplateData) error
	SendOrderPaidNotification(ctx context.Context, templateData *domain.OrderTemplateData) error
}

type OrderPaidConsumerService interface {
	RunOrderConsumer(ctx context.Context) error
}

type ShipAssemblyConsumerService interface {
	RunAssemblyConsumer(ctx context.Context) error
}

type BotService interface{
	Start(ctx context.Context)
}