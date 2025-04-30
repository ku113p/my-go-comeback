package handlers

import (
	"context"
	"crypto/platform/app"
	"crypto/platform/telegram/middleware"
	"crypto/platform/telegram/services"
	"crypto/platform/telegram/view"
	"crypto/platform/utils"
	"errors"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func AttachCallbackQueryData(cb CallbackQueryData, opts []bot.Option, adapter func(HandlerFunc) bot.HandlerFunc) []bot.Option {
	wrappedHandler := adapter(cb.Handle)
	opts = append(opts, bot.WithCallbackQueryDataHandler(cb.Prefix(), bot.MatchTypePrefix, wrappedHandler))
	return opts
}

type NotificationInfoCallbackQueryData struct {
	*app.App
}

func NewNotificationInfoCallbackQueryData(a *app.App) *NotificationInfoCallbackQueryData {
	return &NotificationInfoCallbackQueryData{App: a}
}

func (c *NotificationInfoCallbackQueryData) Prefix() string {
	return "n_"
}

func (c *NotificationInfoCallbackQueryData) Handle(ctx context.Context, b *bot.Bot, update *models.Update, h *middleware.TelegramRequestHelper) {
	h.AnswerCallbackQuery(ctx, update.CallbackQuery.ID)

	s := strings.Replace(update.CallbackQuery.Data, "n_", "", 1)
	s = strings.Trim(s, " ")

	n, err := h.NotificationService.GetNotificationByID(s)
	if err != nil {
		var expectedError *services.ExpectedError
		if errors.As(err, &expectedError) {
			h.SendError(ctx, expectedError.Message)
			return
		}
		h.SendUnexpectedError(ctx, "failed get notification by id", err)
		return
	}

	text := fmt.Sprintf("Notification\nSymbol: %v\nWhen: %v\nAmount: $%v", n.Symbol, n.Sign.When(), utils.FloatComma(n.Amount))
	kb := view.BuildNotificationInfoKeyboard(n)
	h.SendMessageWithMarkup(ctx, text, kb)
}

type RequestDeleteNotificationCallbackQueryData struct {
	*app.App
}

func NewRequestDeleteNotificationCallbackQueryData(a *app.App) *RequestDeleteNotificationCallbackQueryData {
	return &RequestDeleteNotificationCallbackQueryData{App: a}
}

func (c *RequestDeleteNotificationCallbackQueryData) Prefix() string {
	return "rdn_"
}

func (c *RequestDeleteNotificationCallbackQueryData) Handle(ctx context.Context, b *bot.Bot, update *models.Update, h *middleware.TelegramRequestHelper) {
	h.AnswerCallbackQuery(ctx, update.CallbackQuery.ID)

	s := strings.Replace(update.CallbackQuery.Data, "rdn_", "", 1)
	s = strings.Trim(s, " ")

	n, err := h.NotificationService.GetNotificationByID(s)
	if err != nil {
		var expectedError *services.ExpectedError
		if errors.As(err, &expectedError) {
			h.SendError(ctx, expectedError.Message)
			return
		}
		h.SendUnexpectedError(ctx, "failed get notification by id", err)
		return
	}

	text := fmt.Sprintf("Are you sure you want to delete this %v notification?", n.Symbol)
	kb := view.BuildConfirmDeleteNotificationKeyboard(n)
	h.SendMessageWithMarkup(ctx, text, kb)
}

type DeleteNotificationCallbackQueryData struct {
	*app.App
}

func NewDeleteNotificationCallbackQueryData(a *app.App) *DeleteNotificationCallbackQueryData {
	return &DeleteNotificationCallbackQueryData{App: a}
}

func (c *DeleteNotificationCallbackQueryData) Prefix() string {
	return "dn_"
}

func (c *DeleteNotificationCallbackQueryData) Handle(ctx context.Context, b *bot.Bot, update *models.Update, h *middleware.TelegramRequestHelper) {
	h.AnswerCallbackQuery(ctx, update.CallbackQuery.ID)

	s := strings.Replace(update.CallbackQuery.Data, "dn_", "", 1)
	s = strings.Trim(s, " ")

	if err := h.NotificationService.DeleteNotificationByID(s); err != nil {
		var expectedError *services.ExpectedError
		if errors.As(err, &expectedError) {
			h.SendError(ctx, expectedError.Message)
			return
		}
		h.SendUnexpectedError(ctx, "failed delete notification", err)
		return
	}

	h.SendMessage(ctx, "Notification deleted")
}

type DeleteMessageCallbackQueryData struct {
	*app.App
}

func NewDeleteMessageCallbackQueryData(a *app.App) *DeleteMessageCallbackQueryData {
	return &DeleteMessageCallbackQueryData{App: a}
}

func (c *DeleteMessageCallbackQueryData) Prefix() string {
	return "dm_"
}

func (c *DeleteMessageCallbackQueryData) Handle(ctx context.Context, b *bot.Bot, update *models.Update, h *middleware.TelegramRequestHelper) {
	h.AnswerCallbackQuery(ctx, update.CallbackQuery.ID)

	h.DeleteMessage(ctx, update.CallbackQuery.Message.Message.ID)
	h.SendMessage(ctx, "Cancelled")
}
