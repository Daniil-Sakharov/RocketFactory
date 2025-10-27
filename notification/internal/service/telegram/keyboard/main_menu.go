package keyboard

import "github.com/go-telegram/bot/models"

var _ Builder = (*MainMenuBuilder)(nil)

type MainMenuBuilder struct{}

func NewMainMenuBuilder() *MainMenuBuilder {
	return &MainMenuBuilder{}
}

func (b *MainMenuBuilder) BuildMainMenu() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "📦 Мои заказы",
					CallbackData: "my_orders",
				},
				{
					Text:         "ℹ️ Помощь",
					CallbackData: "help",
				},
			},
			{
				{
					Text:         "📊 Статистика",
					CallbackData: "stats",
				},
			},
		},
	}
}

func (b *MainMenuBuilder) BuildOrdersMenu() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "📋 Активные заказы",
					CallbackData: "active_orders",
				},
			},
			{
				{
					Text:         "✅ Завершённые заказы",
					CallbackData: "completed_orders",
				},
			},
			{
				{
					Text:         "🔙 Назад",
					CallbackData: "back_to_main",
				},
			},
		},
	}
}

func (b *MainMenuBuilder) BuildBackButton() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "🔙 Назад в меню",
					CallbackData: "back_to_main",
				},
			},
		},
	}
}
