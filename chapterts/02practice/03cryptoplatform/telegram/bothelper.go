package telegram

import (
	"context"
	"crypto/platform/app"
	"crypto/platform/db"
	"crypto/platform/models"
	"errors"
	"fmt"

	"github.com/go-telegram/bot"
	telegramModels "github.com/go-telegram/bot/models"
)

type BotHelper struct {
	app  *app.App
	mode mode
}

func (h *BotHelper) Run() error {
	opts := []bot.Option{
		bot.WithMiddlewares(h.withUser),
		bot.WithDefaultHandler(h.defaultHandler()),
		bot.WithMessageTextHandler("help", bot.MatchTypeCommand, command),
	}

	return h.mode.runBot(opts...)
}

func (h *BotHelper) defaultHandler() func(context.Context, *bot.Bot, *telegramModels.Update) {
	return func(ctx context.Context, b *bot.Bot, update *telegramModels.Update) {
		u := user(ctx)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("%#v", *u),
		})
	}
}

func command(ctx context.Context, b *bot.Bot, update *telegramModels.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "This bot help to monitor crypto prices",
	})
}

func (h *BotHelper) withUser(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *telegramModels.Update) {
		withUser(ctx, b, update, next, h.app)
	}
}

type userKeyType string

const userKey userKeyType = "user"

func withUser(ctx context.Context, b *bot.Bot, update *telegramModels.Update, next bot.HandlerFunc, a *app.App) {
	if telegramUser := update.Message.From; telegramUser != nil {
		u, err := a.DB.GetUserByTelegramChatID(telegramUser.ID)
		if err != nil {
			if errors.Is(err, db.ErrNotExists) {
				u = models.NewUser(telegramUser.ID)
				u, err = a.DB.CreateUser(u)
				if err != nil {
					a.Logger.Error("failed create user", "error", err)
					return
				}
				a.Logger.Info("created user", "user", u)
			} else {
				a.Logger.Error("failed get user from db", "error", err)
				return
			}
		}
		ctx = context.WithValue(ctx, userKey, u)
	}
	next(ctx, b, update)
}

func user(ctx context.Context) *models.User {
	return ctx.Value(userKey).(*models.User)
}
