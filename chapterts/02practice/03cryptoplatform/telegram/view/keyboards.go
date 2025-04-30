package view

import (
	"crypto/platform/models"
	"fmt"

	telegramModels "github.com/go-telegram/bot/models"
)

func NotificationsKeyboard(ns []*models.Notification) *telegramModels.InlineKeyboardMarkup {
	buttons := make([][]telegramModels.InlineKeyboardButton, 0)
	for _, n := range ns {
		row := []telegramModels.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("%v %s $%v", n.Symbol, n.Sign, n.Amount),
				CallbackData: fmt.Sprintf("n_%v", n.ID.String()),
			},
		}
		buttons = append(buttons, row)
	}

	return &telegramModels.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}
