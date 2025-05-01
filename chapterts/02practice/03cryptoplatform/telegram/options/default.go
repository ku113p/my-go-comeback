package options

import (
	"context"
	"crypto/platform/telegram/handlers"
	"crypto/platform/telegram/helpers"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type defaultParams struct{}

func (p *defaultParams) ToOption(adapter handlers.HandlerAdatper) bot.Option {
	return bot.WithDefaultHandler(adapter(defaultEcho))
}

func NewDefaultParams() OptionParams {
	return &defaultParams{}
}

func defaultEcho(ctx context.Context, update *models.Update, h *helpers.TelegramRequestHelper) {
	if update.Message != nil {
		h.SendMessage(ctx, fmt.Sprintf("%#v", h.User))
	}
}
