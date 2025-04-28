package telegram

import (
	"context"
	"crypto/platform/app"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type BotHelper struct {
	app  *app.App
	mode mode
}

func (h *BotHelper) Run() error {
	opts := []bot.Option{
		bot.WithDefaultHandler(h.defaultHandler()),
		bot.WithMessageTextHandler("foo", bot.MatchTypeCommand, command),
	}

	return h.mode.runBot(opts...)
}

func (h *BotHelper) defaultHandler() func(context.Context, *bot.Bot, *models.Update) {
	return func(ctx context.Context, bot *bot.Bot, update *models.Update) {
		bHandler(ctx, bot, update, h.app)
	}
}

func bHandler(ctx context.Context, b *bot.Bot, update *models.Update, a *app.App) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("%#v", *a),
	})
}

func command(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "bar",
	})
}
