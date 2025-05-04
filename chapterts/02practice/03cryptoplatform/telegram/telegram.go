package telegram

import (
	"context"
	"crypto/platform/app"
	"crypto/platform/models"
	"crypto/platform/telegram/handlers"
	"crypto/platform/telegram/helpers"
	"crypto/platform/telegram/options"
	"crypto/platform/utils"
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

	opts := buildOptions(h.App)

	token, err := token()
	if err != nil {
		return err
	}

	var mode mode
	if useWebhook() {
		mode = modeWebhook
	} else {
		mode = modePooling
	}
	h.Logger.Info(fmt.Sprintf("Telegram bot will be runned in mode %v", mode))

	myBot := newMyBotBuilder().
		withMode(mode).
		withOptions(opts).
		withToken(*token).
		build()

	return myBot.run(ctx)
}

func buildOptions(app *app.App) []bot.Option {
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

func token() (*string, error) {
	token, ok := os.LookupEnv(tokenEnvKey)
	if !ok {
		return nil, fmt.Errorf("env `%s` not found", tokenEnvKey)
	}

	return &token, nil
}

func SendNotification(ctx context.Context, n *models.Notification, app *app.App) error {
	token, err := token()
	if err != nil {
		return err
	}

	bot, err := bot.New(*token)
	if err != nil {
		return err
	}

	user, err := app.DB.GetUserByID(*n.UserID)
	if err != nil {
		return err
	}

	helper := helpers.NewTelegramRequestHelper(bot, user, app)
	text := fmt.Sprintf(
		"Signal!\n%v %v $%v",
		n.Symbol,
		n.Sign.When(),
		utils.FloatComma(n.Amount),
	)
	helper.SendMessage(ctx, text)

	return nil
}
