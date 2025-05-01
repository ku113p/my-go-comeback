package telegram

import (
	"context"

	"github.com/go-telegram/bot"
)

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
