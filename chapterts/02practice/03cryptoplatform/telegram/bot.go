package telegram

import (
	"context"
	"crypto/platform/app"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const tokenEnvKey = "TG_API_TOKEN"
const webhookURLEnvKey = "TG_WEBHOOK_URL"
const webhookPortEnvKey = "TG_WEBHOOK_PORT"

type mode int

const (
	ModePooling mode = iota
	ModeWebhook
)

func (m mode) NewBot(a *app.App) *Bot {
	return &Bot{a, m}
}

func (m mode) runBot(ctx context.Context, b *bot.Bot) error {
	switch m {
	case ModePooling:
		runPooling(ctx, b)
	case ModeWebhook:
		if err := runWebhook(ctx, b); err != nil {
			return err
		}
	}

	return nil
}

func runPooling(ctx context.Context, b *bot.Bot) {
	b.Start(ctx)
}

func runWebhook(ctx context.Context, b *bot.Bot) error {
	url, port, err := getWebhookUrlAndPort()
	if err != nil {
		return err
	}

	b.SetWebhook(ctx, &bot.SetWebhookParams{
		URL: *url,
	})

	go func() {
		http.ListenAndServe(fmt.Sprintf(":%d", *port), b.WebhookHandler())
	}()

	b.StartWebhook(ctx)

	return nil
}

func getWebhookUrlAndPort() (*string, *int, error) {
	url, ok := os.LookupEnv(webhookURLEnvKey)
	if !ok {
		return nil, nil, fmt.Errorf("env `%s` not found", webhookURLEnvKey)
	}

	portStr, ok := os.LookupEnv(webhookPortEnvKey)
	if !ok {
		return nil, nil, fmt.Errorf("env `%s` not found", webhookPortEnvKey)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, nil, err
	}

	return &url, &port, nil
}

type Bot struct {
	app  *app.App
	mode mode
}

func (b *Bot) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
		bot.WithMessageTextHandler("foo", bot.MatchTypeCommand, command),
	}

	botToken, ok := os.LookupEnv(tokenEnvKey)
	if !ok {
		return fmt.Errorf("env `%s` not found", tokenEnvKey)
	}
	bot, err := bot.New(botToken, opts...)
	if err != nil {
		return err
	}

	if err := b.mode.runBot(ctx, bot); err != nil {
		return err
	}

	return nil
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   update.Message.Text,
	})
}

func command(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "bar",
	})
}
