package telegram

import (
	"context"
	"crypto/platform/app"
	"crypto/platform/db"
	"crypto/platform/models"
	"crypto/platform/telegram/middleware"
	"crypto/platform/utils"
	"errors"
	"fmt"
	"strconv"
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
		bot.WithMessageTextHandler("help", bot.MatchTypeCommand, helpCommand),
		bot.WithMessageTextHandler("add", bot.MatchTypeCommand, h.addCommand),
		bot.WithMessageTextHandler("list", bot.MatchTypeCommand, h.listCommand),
		bot.WithCallbackQueryDataHandler("n_", bot.MatchTypePrefix, h.notificationCallbackHandler),
		bot.WithCallbackQueryDataHandler("rdn_", bot.MatchTypePrefix, h.requestDeleteNotificationCallbackHandler),
		bot.WithCallbackQueryDataHandler("dn_", bot.MatchTypePrefix, h.deleteNotificationCallbackHandler),
		bot.WithCallbackQueryDataHandler("dm_", bot.MatchTypePrefix, h.deleteMessageCallbackHandler),
	}

	return h.mode.runBot(opts...)
}

func (h *BotHelper) defaultHandler(ctx context.Context, b *bot.Bot, update *telegramModels.Update) {
	if update.Message != nil {
		telegramHelper := middleware.ContextTelegramRequestHelper(ctx)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("%#v", telegramHelper.User),
		})
	}
}

func helpCommand(ctx context.Context, b *bot.Bot, update *telegramModels.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "This bot help to monitor crypto prices",
	})
}

func (h *BotHelper) withTelegramRequestHelper(next bot.HandlerFunc) bot.HandlerFunc {
	return middleware.WithTelegramRequestHelper(next, h.App)
}

func (h *BotHelper) addCommand(ctx context.Context, b *bot.Bot, update *telegramModels.Update) {
	telegramHelper := middleware.ContextTelegramRequestHelper(ctx)

	s := strings.Replace(update.Message.Text, "/add ", "", 1)
	s = strings.Trim(s, " ")

	n, err := newNotificationFromString(s)
	if err != nil {
		telegramHelper.SendError(fmt.Sprintf("%s", err))
		return
	}

	token, err := h.DB.GetPrice(n.Symbol)
	if err != nil {
		telegramHelper.SendUnexpectedError("failed get price", err)
		return
	}

	if n.Check(token) {
		telegramHelper.SendError("price already reached target amount")
		return
	}

	n.UserID = telegramHelper.User.ID
	n, err = h.DB.CreateNotification(n)
	if err != nil {
		telegramHelper.SendUnexpectedError("failed create notification", err)
		return
	}

	telegramHelper.SendMessage(fmt.Sprintf("Notification #{%s} created.", *n.ID))
}

func newNotificationFromString(s string) (*models.Notification, error) {
	words := strings.SplitN(s, " ", 3)
	if len(words) != 3 {
		return nil, fmt.Errorf("invalid format")
	}

	symbol, signString, amountString := words[0], words[1], words[2]
	symbol = strings.ToUpper(symbol)

	sign, err := models.ParseSign(signString)
	if err != nil {
		return nil, fmt.Errorf("invalid sign")
	}

	amountString = strings.ReplaceAll(amountString, ",", ".")
	amount, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid amount")
	}

	msg := fmt.Sprintf("price %v %v %v", symbol, sign.String(), amount)
	n := models.NewNotification(symbol, *sign, amount, nil, &msg)

	return n, nil
}

func (h *BotHelper) listCommand(ctx context.Context, b *bot.Bot, update *telegramModels.Update) {
	telegramHelper := middleware.ContextTelegramRequestHelper(ctx)

	ns, err := h.DB.ListNotificationsByUserID(*telegramHelper.User.ID)
	if err != nil {
		telegramHelper.SendUnexpectedError("failed list notifications", err)
		return
	}
	kb := notificationsKeyboard(ns)

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        fmt.Sprintf("You have %d Notificatins", len(ns)),
		ReplyMarkup: kb,
	})
	if err != nil {
		h.Logger.Error("failed send message", "error", err)
	}
}

func notificationsKeyboard(ns []*models.Notification) *telegramModels.InlineKeyboardMarkup {
	buttons := make([][]telegramModels.InlineKeyboardButton, 0)
	for _, n := range ns {
		row := []telegramModels.InlineKeyboardButton{
			{
				Text:         fmt.Sprintf("%v %s $%v", n.Symbol, n.Sign, n.Amount),
				CallbackData: fmt.Sprintf("n_%v", n.ID.String()),
			},
		}
		buttons = append(buttons, row)
	}

	return &telegramModels.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
}

func (h *BotHelper) notificationCallbackHandler(ctx context.Context, b *bot.Bot, update *telegramModels.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	telegramHelper := middleware.ContextTelegramRequestHelper(ctx)

	s := strings.Replace(update.CallbackQuery.Data, "n_", "", 1)
	s = strings.Trim(s, " ")

	notificationID, err := uuid.Parse(s)
	if err != nil {
		telegramHelper.SendUnexpectedError("failed parse notification id", err)
		return
	}

	n, err := h.DB.GetNotificationByID(notificationID)
	if err != nil {
		if errors.Is(err, db.ErrNotExists) {
			telegramHelper.SendMessage("notification not found")
			return
		}
		telegramHelper.SendUnexpectedError("failed get notification by id", err)
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
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	telegramHelper := middleware.ContextTelegramRequestHelper(ctx)

	s := strings.Replace(update.CallbackQuery.Data, "rdn_", "", 1)
	s = strings.Trim(s, " ")

	notificationID, err := uuid.Parse(s)
	if err != nil {
		telegramHelper.SendUnexpectedError("failed parse notification id", err)
		return
	}

	n, err := h.DB.GetNotificationByID(notificationID)
	if err != nil {
		if errors.Is(err, db.ErrNotExists) {
			telegramHelper.SendMessage("notification not found")
			return
		}
		telegramHelper.SendUnexpectedError("failed get notification by id", err)
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
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	telegramHelper := middleware.ContextTelegramRequestHelper(ctx)

	s := strings.Replace(update.CallbackQuery.Data, "dn_", "", 1)
	s = strings.Trim(s, " ")

	notificationID, err := uuid.Parse(s)
	if err != nil {
		telegramHelper.SendUnexpectedError("failed parse notification id", err)
		return
	}

	err = h.DB.RemoveNotification(notificationID)
	if err != nil {
		if errors.Is(err, db.ErrNotExists) {
			telegramHelper.SendMessage("notification not found")
			return
		}
		telegramHelper.SendUnexpectedError("failed delete notification", err)
		return
	}

	telegramHelper.SendMessage("Notification deleted")
}

func (h *BotHelper) deleteMessageCallbackHandler(ctx context.Context, b *bot.Bot, update *telegramModels.Update) {
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})

	telegramHelper := middleware.ContextTelegramRequestHelper(ctx)

	if _, err := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
		MessageID: update.CallbackQuery.Message.Message.ID,
	}); err != nil {
		telegramHelper.SendUnexpectedError("failed delete message", err)
		return
	}
	telegramHelper.SendMessage("Cancelled")
}
