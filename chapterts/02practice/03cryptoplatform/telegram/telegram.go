package telegram

import (
	"context"
	"crypto/platform/app"
	"fmt"
	"net/http"
	"os"
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

func (m mode) NewBotRunner(a *app.App) *BotRunner {
	return &BotRunner{a, m}
}

func (m mode) runBot(opts ...bot.Option) error {
	ctx, cancel := context.WithCancel(context.Background())
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
	c, err := getWebhookConnect()
	if err != nil {
		return err
	}

	b.SetWebhook(ctx, &bot.SetWebhookParams{
		URL: c.url,
	})

	go func() {
		http.ListenAndServe(fmt.Sprintf(":%d", c.port), b.WebhookHandler())
	}()

	b.StartWebhook(ctx)

	return nil
}

type webhookConnect struct {
	url  string
	port int
}

func getWebhookConnect() (*webhookConnect, error) {
	url, ok := os.LookupEnv(webhookURLEnvKey)
	if !ok {
		return nil, fmt.Errorf("env `%s` not found", webhookURLEnvKey)
	}

	portStr, ok := os.LookupEnv(webhookPortEnvKey)
	if !ok {
		return nil, fmt.Errorf("env `%s` not found", webhookPortEnvKey)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	return &webhookConnect{url, port}, nil
}
