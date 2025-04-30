package telegram

import (
	"context"
	"crypto/platform/app"
	"crypto/platform/telegram/handlers"
	"crypto/platform/telegram/middleware"
	"fmt"

	"github.com/go-telegram/bot"
	telegramModels "github.com/go-telegram/bot/models"
)

type BotHelper struct {
	*app.App
	mode mode
}

func (h *BotHelper) Run() error {
	opts := []bot.Option{
		bot.WithMiddlewares(h.withTelegramRequestHelper),
		bot.WithDefaultHandler(h.defaultHandler),
	}

	commandHandlers := h.commandHandlers()
	for _, c := range commandHandlers {
		opts = handlers.AttachCommand(c, opts, h.wrapHandler)
	}

	callbackQueryDataHandlers := h.callbackQueryDataHandlers()
	for _, c := range callbackQueryDataHandlers {
		opts = handlers.AttachCallbackQueryData(c, opts, h.wrapHandler)
	}

	return h.mode.runBot(opts...)
}

func (h *BotHelper) commandHandlers() []handlers.Command {
	return []handlers.Command{
		handlers.NewHelpCommand(h.App),
		handlers.NewAddCommand(h.App),
		handlers.NewListCommand(h.App),
	}
}

func (h *BotHelper) callbackQueryDataHandlers() []handlers.CallbackQueryData {
	return []handlers.CallbackQueryData{
		handlers.NewNotificationInfoCallbackQueryData(h.App),
		handlers.NewRequestDeleteNotificationCallbackQueryData(h.App),
		handlers.NewDeleteNotificationCallbackQueryData(h.App),
		handlers.NewDeleteMessageCallbackQueryData(h.App),
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
