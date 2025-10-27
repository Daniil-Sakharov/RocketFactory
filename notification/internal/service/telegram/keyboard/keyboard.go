package keyboard

import "github.com/go-telegram/bot/models"

type Builder interface {
	BuildMainMenu() *models.InlineKeyboardMarkup
	BuildOrdersMenu() *models.InlineKeyboardMarkup
	BuildBackButton() *models.InlineKeyboardMarkup
}
