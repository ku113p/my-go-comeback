package options

import (
	"crypto/platform/app"
	"crypto/platform/telegram/handlers"
	"crypto/platform/telegram/middleware"

	"github.com/go-telegram/bot"
)

type withUserParams struct {
	*app.App
}

func (p *withUserParams) ToOption(adapter handlers.HandlerAdatper) bot.Option {
	return bot.WithMiddlewares(p.withUser)
}

func newWithUserParams(app *app.App) OptionParams {
	return &withUserParams{app}
}

func GetWithUserParamsCreator(app *app.App) OptionParamsBuilder {
	return func() OptionParams {
		return newWithUserParams(app)
	}
}

func (p *withUserParams) withUser(next bot.HandlerFunc) bot.HandlerFunc {
	return middleware.ContextWithUser(next, p.App)
}
