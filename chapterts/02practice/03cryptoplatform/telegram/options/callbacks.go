package options

import (
	"context"
	"crypto/platform/telegram/handlers"
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

type callbackQueryParams struct {
	prefix string
	call   handlers.HandlerFunc
}

func (params *callbackQueryParams) ToOption(adapter handlers.HandlerAdatper) bot.Option {
	return bot.WithCallbackQueryDataHandler(
		params.prefix, bot.MatchTypePrefix, adapter(params.call),
	)
}

type callbackQueryParamsBuilder struct {
	prefix *string
	call   handlers.HandlerFunc
}

func newCallbackQueryParamsBuilder() *callbackQueryParamsBuilder {
	return &callbackQueryParamsBuilder{}
}

func (builder *callbackQueryParamsBuilder) withPrefix(prefix string) *callbackQueryParamsBuilder {
	builder.prefix = &prefix
	return builder
}

func (builder *callbackQueryParamsBuilder) withCall(call handlers.HandlerFunc) *callbackQueryParamsBuilder {
	builder.call = call
	return builder
}

func (builder *callbackQueryParamsBuilder) build() OptionParams {
	return &callbackQueryParams{
		prefix: *builder.prefix,
		call:   builder.call,
	}
}

func NewNotificationInfoCallbackQueryParams() OptionParams {
	return newCallbackQueryParamsBuilder().
		withPrefix("n_").
		withCall(notificationInfo).
		build()
}

func notificationInfo(ctx context.Context, update *models.Update, h *middleware.TelegramRequestHelper) {
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

func NewRequestDeleteNotificationCallbackQueryParams() OptionParams {
	return newCallbackQueryParamsBuilder().
		withPrefix("rdn_").
		withCall(requestDeleteNotification).
		build()
}

func requestDeleteNotification(ctx context.Context, update *models.Update, h *middleware.TelegramRequestHelper) {
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

func NewDeleteNotificationCallbackQueryParams() OptionParams {
	return newCallbackQueryParamsBuilder().
		withPrefix("dn_").
		withCall(deleteNotification).
		build()
}

func deleteNotification(ctx context.Context, update *models.Update, h *middleware.TelegramRequestHelper) {
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

func NewDeleteMessageCallbackQueryParams() OptionParams {
	return newCallbackQueryParamsBuilder().
		withPrefix("dm_").
		withCall(deleteMessage).
		build()
}

func deleteMessage(ctx context.Context, update *models.Update, h *middleware.TelegramRequestHelper) {
	h.AnswerCallbackQuery(ctx, update.CallbackQuery.ID)

	h.DeleteMessage(ctx, update.CallbackQuery.Message.Message.ID)
	h.SendMessage(ctx, "Cancelled")
}
