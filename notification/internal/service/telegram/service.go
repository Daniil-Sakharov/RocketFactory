package telegram

import (
	"context"

	"github.com/Daniil-Sakharov/RocketFactory/notification/internal/client/http"
	"github.com/Daniil-Sakharov/RocketFactory/notification/internal/model/domain"
)

const chatID = 6871748022

type service struct {
	telegramClient http.TelegramClient
	templateEngine *TemplateEngine
}

func NewService(telegramClient http.TelegramClient, templateEngine *TemplateEngine) *service {
	return &service{
		telegramClient: telegramClient,
		templateEngine: templateEngine,
	}
}

func (s *service) SendShipAssembledNotification(ctx context.Context, templateData *domain.AssembledTemplateData) error {
	message, err := s.templateEngine.Render("assembled_notification.tmpl", templateData)
	if err != nil {
		return err
	}

	return s.telegramClient.SendMessage(ctx, chatID, message)
}

func (s *service) SendOrderPaidNotification(ctx context.Context, templateData *domain.OrderTemplateData) error {
	message, err := s.templateEngine.Render("paid_notification.tmpl", templateData)
	if err != nil {
		return err
	}

	return s.telegramClient.SendMessage(ctx, chatID, message)
}
