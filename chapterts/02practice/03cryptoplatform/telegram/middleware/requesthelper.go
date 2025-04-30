package middleware

import (
	"context"
	"crypto/platform/app"
	"crypto/platform/db"
	"crypto/platform/models"
	"crypto/platform/telegram/services"
	"errors"

	"github.com/go-telegram/bot"
	telegramModels "github.com/go-telegram/bot/models"
)

type TelegramRequestHelper struct {
	*app.App
	b                   *bot.Bot
	chatID              int64
	User                *models.User
	NotificationService *services.NotificationService
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
	notificationService := services.NewNotificationService(a)

	return &TelegramRequestHelper{a, b, chatID, u, notificationService}, nil
}

func (h *TelegramRequestHelper) SendMessage(ctx context.Context, text string) {
	_, err := h.b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: h.chatID,
		Text:   text,
	})

	h.logErrorIfNeed(err)
}

func (h *TelegramRequestHelper) logErrorIfNeed(err error) {
	if err != nil {
		h.Logger.Error("failed send message", "error", err)
	}
}

func (h *TelegramRequestHelper) SendError(ctx context.Context, message string) {
	h.SendMessage(ctx, message)
}

func (h *TelegramRequestHelper) SendUnexpectedError(ctx context.Context, subject string, err error) {
	h.Logger.Error(subject, "error", err)
	h.SendMessage(ctx, "Unexpected error occurred")
}

func (h *TelegramRequestHelper) AnswerCallbackQuery(ctx context.Context, callbackQueryID string) {
	_, err := h.b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQueryID,
		ShowAlert:       false,
	})

	h.logErrorIfNeed(err)
}

func (h *TelegramRequestHelper) SendMessageWithMarkup(ctx context.Context, text string, kb telegramModels.ReplyMarkup) {
	_, err := h.b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      h.chatID,
		Text:        text,
		ReplyMarkup: kb,
	})

	h.logErrorIfNeed(err)
}

func (h *TelegramRequestHelper) DeleteMessage(ctx context.Context, messageID int) {
	_, err := h.b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    h.chatID,
		MessageID: messageID,
	})

	h.logErrorIfNeed(err)
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
