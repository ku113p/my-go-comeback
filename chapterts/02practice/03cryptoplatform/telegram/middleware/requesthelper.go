package middleware

import (
	"context"
	"crypto/platform/app"
	"crypto/platform/db"
	"crypto/platform/models"
	"errors"
	"log/slog"

	"github.com/go-telegram/bot"
	telegramModels "github.com/go-telegram/bot/models"
)

type TelegramRequestHelper struct {
	*app.App
	b      *bot.Bot
	chatID int64
	User   *models.User
}

func newTelegramRequestHelper(b *bot.Bot, chatID int64, a *app.App) (*TelegramRequestHelper, error) {
	u, err := a.DB.GetUserByTelegramChatID(chatID)
	if err != nil {
		if errors.Is(err, db.ErrNotExists) {
			u = models.NewUser(chatID)
			u, err = a.DB.CreateUser(u)
			if err != nil {
				return nil, err
			}
			a.Logger.Info("created user", "user", u)
		} else {
			return nil, err
		}
	}

	return &TelegramRequestHelper{a, b, chatID, u}, nil
}

func (h *TelegramRequestHelper) SendMessage(ctx context.Context, text string) {
	sendMessageFunc(ctx, h.chatID, h.b, h.Logger)(text)
}

func (h *TelegramRequestHelper) SendError(ctx context.Context, message string) {
	h.SendMessage(ctx, message)
}

func (h *TelegramRequestHelper) SendUnexpectedError(ctx context.Context, subject string, err error) {
	h.Logger.Error(subject, "error", err)
	h.SendMessage(ctx, "Unexpected error occurred")
}

func sendMessageFunc(ctx context.Context, chatID int64, b *bot.Bot, logger *slog.Logger) func(string) {
	return func(msg string) {
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   msg,
		})

		if err != nil {
			logger.Error("failed send message", "error", err)
		}
	}
}

func (h *TelegramRequestHelper) AnswerCallbackQuery(ctx context.Context, callbackQueryID string) (bool, error) {
	return h.b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQueryID,
		ShowAlert:       false,
	})
}

type helperKeyType string

const helperKey helperKeyType = "telegram helper"

func ContextTelegramRequestHelper(ctx context.Context) *TelegramRequestHelper {
	value := ctx.Value(helperKey)
	h, _ := value.(*TelegramRequestHelper)
	return h
}

func WithTelegramRequestHelper(next bot.HandlerFunc, a *app.App) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *telegramModels.Update) {
		var chatID int64

		switch {
		case update.Message != nil:
			chatID = update.Message.Chat.ID
		case update.CallbackQuery != nil:
			chatID = update.CallbackQuery.Message.Message.Chat.ID
		default:
			a.Logger.Warn("unable to determine chatID for handler")
			return
		}

		h, err := newTelegramRequestHelper(b, chatID, a)
		if err != nil {
			a.Logger.Error("failed create telegram TelegramRequestHelper", "error", err)
			return
		}
		ctx = context.WithValue(ctx, helperKey, h)

		next(ctx, b, update)
	}
}
