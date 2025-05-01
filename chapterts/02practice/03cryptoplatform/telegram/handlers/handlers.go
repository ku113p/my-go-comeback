package handlers

import (
	"context"
	"crypto/platform/app"
	"crypto/platform/telegram/helpers"
	"crypto/platform/telegram/middleware"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type HandlerFunc func(ctx context.Context, update *models.Update, h *helpers.TelegramRequestHelper)

type HandlerAdatper func(HandlerFunc) bot.HandlerFunc

func GetAdapter(app *app.App) HandlerAdatper {
	return func(fn HandlerFunc) bot.HandlerFunc {
		return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
			user := middleware.ContextUser(ctx)
			telegramHelper := helpers.NewTelegramRequestHelper(bot, user, app)

			fn(ctx, update, telegramHelper)
		}
	}
}
