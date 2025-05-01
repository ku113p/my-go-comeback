package options

import (
	"crypto/platform/telegram/handlers"

	"github.com/go-telegram/bot"
)

type OptionParams interface {
	ToOption(handlers.HandlerAdatper) bot.Option
}

type OptionParamsBuilder func() OptionParams
