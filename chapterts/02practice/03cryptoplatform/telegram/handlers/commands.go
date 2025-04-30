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

type HelpCommand struct {
	*app.App
}

func NewHelpCommand(a *app.App) *HelpCommand {
	return &HelpCommand{a}
}

func (c *HelpCommand) Command() string {
	return "help"
}

func (c *HelpCommand) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	h := middleware.ContextTelegramRequestHelper(ctx)

	h.SendMessage("This bot help to monitor crypto prices")
}

type AddCommand struct {
	*app.App
}

func NewAddCommand(a *app.App) *AddCommand {
	return &AddCommand{a}
}

func (c *AddCommand) Command() string {
	return "add"
}

func (c *AddCommand) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	telegramHelper := middleware.ContextTelegramRequestHelper(ctx)

	s := strings.Replace(update.Message.Text, "/add ", "", 1)
	s = strings.Trim(s, " ")

	n, err := newNotificationFromString(s)
	if err != nil {
		telegramHelper.SendError(fmt.Sprintf("%s", err))
		return
	}

	token, err := c.DB.GetPrice(n.Symbol)
	if err != nil {
		telegramHelper.SendUnexpectedError("failed get price", err)
		return
	}

	if n.Check(token) {
		telegramHelper.SendError("price already reached target amount")
		return
	}

	n.UserID = telegramHelper.User.ID
	n, err = c.DB.CreateNotification(n)
	if err != nil {
		telegramHelper.SendUnexpectedError("failed create notification", err)
		return
	}

	telegramHelper.SendMessage(fmt.Sprintf("Notification #{%s} created.", *n.ID))
}

type ListCommand struct {
	*app.App
}

func NewListCommand(a *app.App) *ListCommand {
	return &ListCommand{a}
}

func (c *ListCommand) Command() string {
	return "list"
}

func (c *ListCommand) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	telegramHelper := middleware.ContextTelegramRequestHelper(ctx)

	ns, err := c.DB.ListNotificationsByUserID(*telegramHelper.User.ID)
	if err != nil {
		telegramHelper.SendUnexpectedError("failed list notifications", err)
		return
	}
	kb := view.NotificationsKeyboard(ns)

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        fmt.Sprintf("You have %d Notificatins", len(ns)),
		ReplyMarkup: kb,
	})
	if err != nil {
		c.Logger.Error("failed send message", "error", err)
	}
}
