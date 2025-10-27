package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/RocketFactory/notification/internal/service/telegram/message"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
)

var _ CallbackHandler = (*ButtonCallbackHandler)(nil)

type ButtonCallbackHandler struct {
	messageBuilder *message.Builder
}

func NewButtonCallbackHandler(messageBuilder *message.Builder) *ButtonCallbackHandler {
	return &ButtonCallbackHandler{
		messageBuilder: messageBuilder,
	}
}

func (h *ButtonCallbackHandler) HandleCallback(ctx context.Context, b *bot.Bot, update *models.Update) {
	callbackQuery := update.CallbackQuery
	callbackData := callbackQuery.Data

	chatID := callbackQuery.From.ID

	logger.Info(ctx, "Callback received",
		zap.Int64("chat_id", chatID),
		zap.String("callback_data", callbackData),
	)

	_, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
	})
	if err != nil {
		logger.Error(ctx, "Failed to answer callback", zap.Error(err))
		return
	}

	var responseText string

	switch callbackData {
	case "my_orders":
		responseText = h.messageBuilder.MyOrdersMessage()
	case "help":
		responseText = h.messageBuilder.HelpMessage()
	case "stats":
		responseText = h.messageBuilder.StatsMessage()
	default:
		responseText = h.messageBuilder.UnknownCommandMessage()
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    chatID,
		Text:      responseText,
		ParseMode: "Markdown",
	})
	if err != nil {
		logger.Error(ctx, "Failed to send callback response", zap.Error(err))
	}
}
