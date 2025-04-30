package handlers

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type CommandHandler interface {
	Command() string
	Handle(ctx context.Context, b *bot.Bot, update *models.Update)
}
