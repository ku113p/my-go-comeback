package handlers

import (
	"context"
	"crypto/platform/telegram/middleware"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type HandlerAdatper func(HandlerFunc) bot.HandlerFunc

type HandlerFunc func(ctx context.Context, update *models.Update, h *middleware.TelegramRequestHelper)

type Handler interface {
	Handle(ctx context.Context, b *bot.Bot, update *models.Update, h *middleware.TelegramRequestHelper)
}
