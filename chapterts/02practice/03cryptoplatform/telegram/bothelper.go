package telegram

import (
	"context"
	"crypto/platform/app"
	"crypto/platform/db"
	"crypto/platform/telegram/handlers"
	"crypto/platform/telegram/middleware"
	"crypto/platform/telegram/view"
	"crypto/platform/utils"
	"errors"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	telegramModels "github.com/go-telegram/bot/models"
	"github.com/google/uuid"
)

type BotHelper struct {
	*app.App
	mode mode
}

func (h *BotHelper) Run() error {
	opts := []bot.Option{
		bot.WithMiddlewares(h.withTelegramRequestHelper),
		bot.WithDefaultHandler(h.defaultHandler),
		bot.WithCallbackQueryDataHandler("n_", bot.MatchTypePrefix, h.notificationCallbackHandler),
		bot.WithCallbackQueryDataHandler("rdn_", bot.MatchTypePrefix, h.requestDeleteNotificationCallbackHandler),
		bot.WithCallbackQueryDataHandler("dn_", bot.MatchTypePrefix, h.deleteNotificationCallbackHandler),
		bot.WithCallbackQueryDataHandler("dm_", bot.MatchTypePrefix, h.deleteMessageCallbackHandler),
	}

	commandHandlers := h.commandHandlers()
	for _, c := range commandHandlers {
		wrappedHandler := h.wrapHandler(c.Handle)
		opts = append(opts, bot.WithMessageTextHandler(c.Command(), bot.MatchTypeCommand, wrappedHandler))
	}

	return h.mode.runBot(opts...)
}

func (h *BotHelper) commandHandlers() []handlers.CommandHandler {
	return []handlers.CommandHandler{
		handlers.NewHelpCommand(h.App),
		handlers.NewAddCommand(h.App),
		handlers.NewListCommand(h.App),
	}
}

func (h *BotHelper) wrapHandler(fn handlers.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *telegramModels.Update) {
		telegramHelper := middleware.ContextTelegramRequestHelper(ctx)
		fn(ctx, b, update, telegramHelper)
	}
}

func (h *BotHelper) defaultHandler(ctx context.Context, b *bot.Bot, update *telegramModels.Update) {
	if update.Message != nil {
		telegramHelper := middleware.ContextTelegramRequestHelper(ctx)
		telegramHelper.SendMessage(ctx, fmt.Sprintf("%#v", telegramHelper.User))
	}
}

func (h *BotHelper) withTelegramRequestHelper(next bot.HandlerFunc) bot.HandlerFunc {
	return middleware.WithTelegramRequestHelper(next, h.App)
}

func (h *BotHelper) notificationCallbackHandler(ctx context.Context, b *bot.Bot, update *telegramModels.Update) {
	telegramHelper := middleware.ContextTelegramRequestHelper(ctx)
	telegramHelper.AnswerCallbackQuery(ctx, update.CallbackQuery.ID)

	s := strings.Replace(update.CallbackQuery.Data, "n_", "", 1)
	s = strings.Trim(s, " ")

	notificationID, err := uuid.Parse(s)
	if err != nil {
		telegramHelper.SendUnexpectedError(ctx, "failed parse notification id", err)
		return
	}

	n, err := h.DB.GetNotificationByID(notificationID)
	if err != nil {
		if errors.Is(err, db.ErrNotExists) {
			telegramHelper.SendMessage(ctx, "notification not found")
			return
		}
		telegramHelper.SendUnexpectedError(ctx, "failed get notification by id", err)
		return
	}

	text := fmt.Sprintf("Notification\nSymbol: %v\nWhen: %v\nAmount: $%v", n.Symbol, n.Sign.When(), utils.FloatComma(n.Amount))
	kb := view.BuildNotificationInfoKeyboard(n)
	telegramHelper.SendMessageWithMarkup(ctx, text, kb)
}

func (h *BotHelper) requestDeleteNotificationCallbackHandler(ctx context.Context, b *bot.Bot, update *telegramModels.Update) {
	telegramHelper := middleware.ContextTelegramRequestHelper(ctx)
	telegramHelper.AnswerCallbackQuery(ctx, update.CallbackQuery.ID)

	s := strings.Replace(update.CallbackQuery.Data, "rdn_", "", 1)
	s = strings.Trim(s, " ")

	notificationID, err := uuid.Parse(s)
	if err != nil {
		telegramHelper.SendUnexpectedError(ctx, "failed parse notification id", err)
		return
	}

	n, err := h.DB.GetNotificationByID(notificationID)
	if err != nil {
		if errors.Is(err, db.ErrNotExists) {
			telegramHelper.SendMessage(ctx, "notification not found")
			return
		}
		telegramHelper.SendUnexpectedError(ctx, "failed get notification by id", err)
		return
	}

	text := fmt.Sprintf("Are you sure you want to delete this %v notification?", n.Symbol)
	kb := view.BuildConfirmDeleteNotificationKeyboard(n)
	telegramHelper.SendMessageWithMarkup(ctx, text, kb)
}

func (h *BotHelper) deleteNotificationCallbackHandler(ctx context.Context, b *bot.Bot, update *telegramModels.Update) {
	telegramHelper := middleware.ContextTelegramRequestHelper(ctx)
	telegramHelper.AnswerCallbackQuery(ctx, update.CallbackQuery.ID)

	s := strings.Replace(update.CallbackQuery.Data, "dn_", "", 1)
	s = strings.Trim(s, " ")

	notificationID, err := uuid.Parse(s)
	if err != nil {
		telegramHelper.SendUnexpectedError(ctx, "failed parse notification id", err)
		return
	}

	err = h.DB.RemoveNotification(notificationID)
	if err != nil {
		if errors.Is(err, db.ErrNotExists) {
			telegramHelper.SendMessage(ctx, "notification not found")
			return
		}
		telegramHelper.SendUnexpectedError(ctx, "failed delete notification", err)
		return
	}

	telegramHelper.SendMessage(ctx, "Notification deleted")
}

func (h *BotHelper) deleteMessageCallbackHandler(ctx context.Context, b *bot.Bot, update *telegramModels.Update) {
	telegramHelper := middleware.ContextTelegramRequestHelper(ctx)
	telegramHelper.AnswerCallbackQuery(ctx, update.CallbackQuery.ID)

	telegramHelper.DeleteMessage(ctx, update.CallbackQuery.Message.Message.ID)
	telegramHelper.SendMessage(ctx, "Cancelled")
}
