package telegram

import (
	"context"

	"github.com/go-telegram/bot"

	serv "github.com/Daniil-Sakharov/RocketFactory/notification/internal/service"
	"github.com/Daniil-Sakharov/RocketFactory/notification/internal/service/telegram/handler"
	"github.com/Daniil-Sakharov/RocketFactory/notification/internal/service/telegram/keyboard"
	"github.com/Daniil-Sakharov/RocketFactory/notification/internal/service/telegram/message"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
)

var _ serv.BotService = (*botService)(nil)

type botService struct {
	bot             *bot.Bot
	startHandler    handler.Handler
	callbackHandler handler.CallbackHandler
}

func NewBotService(bot *bot.Bot, templateEngine *TemplateEngine) *botService {
	keyboardBuilder := keyboard.NewMainMenuBuilder()
	messageBuilder := message.NewBuilder()

	return &botService{
		bot:             bot,
		startHandler:    handler.NewStartHandler(keyboardBuilder, templateEngine),
		callbackHandler: handler.NewButtonCallbackHandler(messageBuilder),
	}
}

func (s *botService) RegisterHandlers() {
	s.bot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, s.startHandler.Handle)
	s.bot.RegisterHandler(bot.HandlerTypeCallbackQueryData, "", bot.MatchTypePrefix, s.callbackHandler.HandleCallback)
}

func (s *botService) Start(ctx context.Context) {
	logger.Info(ctx, "🤖 Starting Telegram bot polling...")

	s.RegisterHandlers()

	s.bot.Start(ctx)

	logger.Info(ctx, "✅ Telegram bot polling started")
}
