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
)

const tokenEnvKey = "TG_API_TOKEN"
const webhookURLEnvKey = "TG_WEBHOOK_URL"
const webhookPortEnvKey = "TG_WEBHOOK_PORT"

type mode int

const (
	ModePooling mode = iota
	ModeWebhook
)

func (m mode) NewBotHelper(a *app.App) *BotHelper {
	return &BotHelper{a, m}
}

func (m mode) runBot(opts ...bot.Option) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	botToken, ok := os.LookupEnv(tokenEnvKey)
	if !ok {
		return fmt.Errorf("env `%s` not found", tokenEnvKey)
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		return err
	}

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
