package telegram

import (
	"context"
	"crypto/platform/app"
	"crypto/platform/telegram/handlers"
	"crypto/platform/telegram/middleware"
	"crypto/platform/telegram/options"
	"fmt"
	"os"

	"github.com/go-telegram/bot"
	telegramModels "github.com/go-telegram/bot/models"
)

const tokenEnvKey = "TG_API_TOKEN"

type BotRunner struct {
	*app.App
}

func NewBotRunner(app *app.App) *BotRunner {
	return &BotRunner{app}
}

func (h *BotRunner) Run() error {
	ctx := context.Background()

	opts := h.options()

	token, err := getToken()
	if err != nil {
		return err
	}

	myBot := newMyBotBuilder().
		withMode(modePooling).
		withOptions(opts).
		withToken(*token).
		build()

	return myBot.run(ctx)
}

func (h *BotRunner) options() []bot.Option {
	opts := []bot.Option{
		bot.WithMiddlewares(h.withTelegramRequestHelper),
		bot.WithDefaultHandler(h.defaultHandler),
	}

	optionsCreators := []options.OptionParamsBuilder{
		options.NewHelpCommandParams,
		options.NewAddCommandParams,
		options.NewListCommandParams,
		options.NewDeleteMessageCallbackQueryParams,
		options.NewDeleteNotificationCallbackQueryParams,
		options.NewNotificationInfoCallbackQueryParams,
		options.NewRequestDeleteNotificationCallbackQueryParams,
	}
	for _, paramCreator := range optionsCreators {
		opts = append(opts, paramCreator().ToOption(h.wrapHandler))
	}

	return opts
}

func (h *BotRunner) wrapHandler(fn handlers.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *telegramModels.Update) {
		telegramHelper := middleware.ContextTelegramRequestHelper(ctx)
		fn(ctx, update, telegramHelper)
	}
}

func (h *BotRunner) defaultHandler(ctx context.Context, b *bot.Bot, update *telegramModels.Update) {
	if update.Message != nil {
		telegramHelper := middleware.ContextTelegramRequestHelper(ctx)
		telegramHelper.SendMessage(ctx, fmt.Sprintf("%#v", telegramHelper.User))
	}
}

func (h *BotRunner) withTelegramRequestHelper(next bot.HandlerFunc) bot.HandlerFunc {
	return middleware.WithTelegramRequestHelper(next, h.App)
}

func getToken() (*string, error) {
	token, ok := os.LookupEnv(tokenEnvKey)
	if !ok {
		return nil, fmt.Errorf("env `%s` not found", tokenEnvKey)
	}

	return &token, nil
}
