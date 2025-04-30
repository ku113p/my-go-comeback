package telegram

import (
	"context"
	"crypto/platform/app"
	"crypto/platform/db"
	"crypto/platform/telegram/handlers"
	"crypto/platform/telegram/middleware"
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
	kb := &telegramModels.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegramModels.InlineKeyboardButton{
			{
				{
					Text:         "Delete ❌",
					CallbackData: fmt.Sprintf("rdn_%v", n.ID.String()),
				},
			},
		},
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        text,
		ReplyMarkup: kb,
	})
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
	kb := &telegramModels.InlineKeyboardMarkup{
		InlineKeyboard: [][]telegramModels.InlineKeyboardButton{
			{
				{
					Text:         "Delete ❌",
					CallbackData: fmt.Sprintf("dn_%v", n.ID.String()),
				},
				{
					Text:         "Cancel ⭕",
					CallbackData: fmt.Sprintf("dm_%v", nil),
				},
			},
		},
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        text,
		ReplyMarkup: kb,
	})
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

	if _, err := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
	}); err != nil {
		telegramHelper.SendUnexpectedError(ctx, "failed delete message", err)
		return
	}
	telegramHelper.SendMessage(ctx, "Cancelled")
}
