package telegram

import (
	"context"
	"crypto/platform/app"
	"crypto/platform/telegram/handlers"
	"crypto/platform/telegram/options"
	"fmt"
	"os"

	"github.com/go-telegram/bot"
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

	opts := getOptions(h.App)

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

func getOptions(app *app.App) []bot.Option {
	adapter := handlers.GetAdapter(app)
	opts := []bot.Option{}

	optionsCreators := []options.OptionParamsBuilder{
		options.GetWithUserParamsCreator(app),
		options.NewDefaultParams,
		options.NewHelpCommandParams,
		options.NewAddCommandParams,
		options.NewListCommandParams,
		options.NewDeleteMessageCallbackQueryParams,
		options.NewDeleteNotificationCallbackQueryParams,
		options.NewNotificationInfoCallbackQueryParams,
		options.NewRequestDeleteNotificationCallbackQueryParams,
	}
	for _, paramCreator := range optionsCreators {
		opts = append(opts, paramCreator().ToOption(adapter))
	}

	return opts
}

func getToken() (*string, error) {
	token, ok := os.LookupEnv(tokenEnvKey)
	if !ok {
		return nil, fmt.Errorf("env `%s` not found", tokenEnvKey)
	}

	return &token, nil
}
