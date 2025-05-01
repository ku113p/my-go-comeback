package helpers

import (
	"context"
	"crypto/platform/app"
	"crypto/platform/models"
	"crypto/platform/telegram/services"

	"github.com/go-telegram/bot"
	telegramModels "github.com/go-telegram/bot/models"
)

type TelegramRequestHelper struct {
	*app.App
	bot      *bot.Bot
	User     *models.User
	Services *services.Services
}

func NewTelegramRequestHelper(bot *bot.Bot, user *models.User, a *app.App) *TelegramRequestHelper {
	notificationService := services.NewServices(a)

	return &TelegramRequestHelper{a, bot, user, notificationService}
}

func (h *TelegramRequestHelper) SendMessage(ctx context.Context, text string) {
	_, err := h.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: *h.User.TelegramChatID,
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
	_, err := h.bot.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQueryID,
		ShowAlert:       false,
	})

	h.logErrorIfNeed(err)
}

func (h *TelegramRequestHelper) SendMessageWithMarkup(ctx context.Context, text string, kb telegramModels.ReplyMarkup) {
	_, err := h.bot.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      *h.User.TelegramChatID,
		Text:        text,
		ReplyMarkup: kb,
	})

	h.logErrorIfNeed(err)
}

func (h *TelegramRequestHelper) DeleteMessage(ctx context.Context, messageID int) {
	_, err := h.bot.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    *h.User.TelegramChatID,
		MessageID: messageID,
	})

	h.logErrorIfNeed(err)
}
