package options

import (
	"context"
	"crypto/platform/telegram/handlers"
	"crypto/platform/telegram/middleware"
	"crypto/platform/telegram/services"
	"crypto/platform/telegram/view"
	"errors"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type commandParams struct {
	name string
	call handlers.HandlerFunc
}

func (params *commandParams) ToOption(adapter handlers.HandlerAdatper) bot.Option {
	return bot.WithMessageTextHandler(
		params.name, bot.MatchTypeCommand, adapter(params.call),
	)
}

type commandParamsBuilder struct {
	name *string
	call handlers.HandlerFunc
}

func newCommandParamsBuilder() *commandParamsBuilder {
	return &commandParamsBuilder{}
}

func (builder *commandParamsBuilder) withName(name string) *commandParamsBuilder {
	builder.name = &name
	return builder
}

func (builder *commandParamsBuilder) withCall(call handlers.HandlerFunc) *commandParamsBuilder {
	builder.call = call
	return builder
}

func (builder *commandParamsBuilder) build() *commandParams {
	return &commandParams{
		name: *builder.name,
		call: builder.call,
	}
}

func NewHelpCommandParams() OptionParams {
	return newCommandParamsBuilder().
		withName("help").
		withCall(help).
		build()
}

func help(ctx context.Context, _ *models.Update, h *middleware.TelegramRequestHelper) {
	h.SendMessage(ctx, "This bot help to monitor crypto prices")
}

func NewAddCommandParams() OptionParams {
	return newCommandParamsBuilder().
		withName("add").
		withCall(add).
		build()
}

func add(ctx context.Context, update *models.Update, h *middleware.TelegramRequestHelper) {
	s := strings.Replace(update.Message.Text, "/add ", "", 1)
	s = strings.Trim(s, " ")

	n, err := h.NotificationService.CreateNotification(h.User, s)
	if err != nil {
		var expectedError *services.ExpectedError
		if errors.As(err, &expectedError) {
			h.SendError(ctx, expectedError.Message)
			return
		}

		h.SendUnexpectedError(ctx, "failed to create notification", err)
		return
	}

	h.SendMessage(ctx, fmt.Sprintf("Notification #{%s} created.", *n.ID))
}

func NewListCommandParams() OptionParams {
	return newCommandParamsBuilder().
		withName("list").
		withCall(list).
		build()
}

func list(ctx context.Context, update *models.Update, h *middleware.TelegramRequestHelper) {
	ns, err := h.NotificationService.GetNotificationsByUser(h.User)
	if err != nil {
		h.SendUnexpectedError(ctx, "failed get list notifications", err)
		return
	}

	text := fmt.Sprintf("You have %d Notificatins", len(ns))
	kb := view.BuildNotificationsKeyboard(ns)
	h.SendMessageWithMarkup(ctx, text, kb)
}
