package telegram

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-telegram/bot"
)

const webhookUrlEnv = "TG_WEBHOOK_URL"

type mode string

const (
	modePooling mode = "polling"
	modeWebhook mode = "webhook"
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
	url, ok := os.LookupEnv(webhookUrlEnv)
	if !ok {
		return fmt.Errorf("env `%s` not found", webhookUrlEnv)
	}

	b.SetWebhook(ctx, &bot.SetWebhookParams{
		URL: url,
	})

	go func() {
		http.ListenAndServe(":8080", b.WebhookHandler())
	}()

	b.StartWebhook(ctx)

	return nil
}

func useWebhook() bool {
	_, ok := os.LookupEnv(webhookUrlEnv)

	return ok
}
