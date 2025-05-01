package telegram

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-telegram/bot"
)

const webhookURLEnvKey = "TG_WEBHOOK_URL"
const webhookPortEnvKey = "TG_WEBHOOK_PORT"

type mode int

const (
	modePooling mode = iota
	modeWebhook
)

func (mode mode) startTelegramBot(ctx context.Context, telegramBot *bot.Bot) error {
	switch mode {
	case modePooling:
		runPooling(ctx, telegramBot)
	case modeWebhook:
		if err := runWebhook(ctx, telegramBot); err != nil {
			return err
		}
	}

	return fmt.Errorf("unknown mode: %v", mode)
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
