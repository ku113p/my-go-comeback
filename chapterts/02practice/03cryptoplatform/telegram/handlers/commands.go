package handlers

import (
	"context"
	"crypto/platform/app"
	"crypto/platform/telegram/middleware"
	"crypto/platform/telegram/view"
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

	n, err := newNotificationFromString(s)
	if err != nil {
		h.SendError(ctx, fmt.Sprintf("%s", err))
		return
	}

	token, err := c.DB.GetPrice(n.Symbol)
	if err != nil {
		h.SendUnexpectedError(ctx, "failed get price", err)
		return
	}

	if n.Check(token) {
		h.SendError(ctx, "price already reached target amount")
		return
	}

	n.UserID = h.User.ID
	n, err = c.DB.CreateNotification(n)
	if err != nil {
		h.SendUnexpectedError(ctx, "failed create notification", err)
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
	ns, err := c.DB.ListNotificationsByUserID(*h.User.ID)
	if err != nil {
		h.SendUnexpectedError(ctx, "failed list notifications", err)
		return
	}
	kb := view.BuildNotificationsKeyboard(ns)

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        fmt.Sprintf("You have %d Notificatins", len(ns)),
		ReplyMarkup: kb,
	})
	if err != nil {
		c.Logger.Error("failed send message", "error", err)
	}
}
