package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/RocketFactory/notification/internal/service/telegram/keyboard"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
)

var _ Handler = (*StartHandler)(nil)

type TemplateRenderer interface {
	Render(templateName string, data interface{}) (string, error)
}

type StartHandler struct {
	keyboardBuilder keyboard.Builder
	templateEngine  TemplateRenderer
}

func NewStartHandler(keyboardBuilder keyboard.Builder, templateEngine TemplateRenderer) *StartHandler {
	return &StartHandler{
		keyboardBuilder: keyboardBuilder,
		templateEngine:  templateEngine,
	}
}

func (h *StartHandler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatID := update.Message.Chat.ID
	userName := update.Message.From.FirstName

	logger.Info(ctx, "User started bot",
		zap.Int64("chat_id", chatID),
		zap.String("username", userName),
	)

	welcomeText, err := h.templateEngine.Render("welcome.tmpl", map[string]string{
		"UserName": userName,
	})
	if err != nil {
		logger.Error(ctx, "Failed to render welcome template", zap.Error(err))
		return
	}

	mainMenu := h.keyboardBuilder.BuildMainMenu()

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID,
		Text:        welcomeText,
		ParseMode:   "Markdown",
		ReplyMarkup: mainMenu,
	})
	if err != nil {
		logger.Error(ctx, "Failed to send welcome message", zap.Error(err))
	}
}
