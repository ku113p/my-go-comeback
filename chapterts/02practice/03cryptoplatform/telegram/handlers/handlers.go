package handlers

import (
	"context"
	"crypto/platform/telegram/middleware"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type HandlerFunc func(ctx context.Context, b *bot.Bot, update *models.Update, h *middleware.TelegramRequestHelper)

type CommandHandler interface {
	Command() string
	Handle(ctx context.Context, b *bot.Bot, update *models.Update, h *middleware.TelegramRequestHelper)
}
