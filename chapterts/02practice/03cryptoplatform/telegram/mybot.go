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

type myBot struct {
	mode  mode
	token string
	opts  []bot.Option
}

func (myBot *myBot) run(ctx context.Context) error {
	telegramBot, err := bot.New(myBot.token, myBot.opts...)
	if err != nil {
		return err
	}

	return myBot.mode.startTelegramBot(ctx, telegramBot)
}

type myBotBuilder struct {
	mode  *mode
	token *string
	opts  []bot.Option
}

func newMyBotBuilder() *myBotBuilder {
	return &myBotBuilder{}
}

func (b *myBotBuilder) withMode(mode mode) *myBotBuilder {
	b.mode = &mode
	return b
}

func (b *myBotBuilder) withOptions(opts []bot.Option) *myBotBuilder {
	b.opts = opts
	return b
}

func (b *myBotBuilder) withToken(token string) *myBotBuilder {
	b.token = &token
	return b
}

func (builder *myBotBuilder) build() *myBot {
	return &myBot{
		mode:  *builder.mode,
		token: *builder.token,
		opts:  builder.opts,
	}
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
