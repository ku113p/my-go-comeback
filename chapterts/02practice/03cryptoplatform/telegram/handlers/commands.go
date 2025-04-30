package handlers

import (
	"context"
	"crypto/platform/app"
	"crypto/platform/telegram/middleware"
	"crypto/platform/telegram/services"
	"crypto/platform/telegram/view"
	"errors"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func AttachCommand(cmd Command, opts []bot.Option, adapter func(HandlerFunc) bot.HandlerFunc) []bot.Option {
	wrappedHandler := adapter(cmd.Handle)
	opts = append(opts, bot.WithMessageTextHandler(cmd.Name(), bot.MatchTypeCommand, wrappedHandler))
	return opts
}

type HelpCommand struct {
	*app.App
}

func NewHelpCommand(a *app.App) *HelpCommand {
	return &HelpCommand{App: a}
}

func (c *HelpCommand) Name() string {
	return "help"
}

func (c *HelpCommand) Handle(ctx context.Context, b *bot.Bot, update *models.Update, h *middleware.TelegramRequestHelper) {
	h.SendMessage(ctx, "This bot help to monitor crypto prices")
}

type AddCommand struct {
	*app.App
}

func NewAddCommand(a *app.App) *AddCommand {
	return &AddCommand{a}
}

func (c *AddCommand) Name() string {
	return "add"
}

func (c *AddCommand) Handle(ctx context.Context, b *bot.Bot, update *models.Update, h *middleware.TelegramRequestHelper) {
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

type ListCommand struct {
	*app.App
}

func NewListCommand(a *app.App) *ListCommand {
	return &ListCommand{a}
}

func (c *ListCommand) Name() string {
	return "list"
}

func (c *ListCommand) Handle(ctx context.Context, b *bot.Bot, update *models.Update, h *middleware.TelegramRequestHelper) {
	ns, err := h.NotificationService.GetNotificationsByUser(h.User)
	if err != nil {
		h.SendUnexpectedError(ctx, "failed get list notifications", err)
		return
	}

	text := fmt.Sprintf("You have %d Notificatins", len(ns))
	kb := view.BuildNotificationsKeyboard(ns)
	h.SendMessageWithMarkup(ctx, text, kb)
}
