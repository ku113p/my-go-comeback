package middleware

import (
	"context"
	"crypto/platform/app"
	"crypto/platform/db"
	"crypto/platform/models"
	"errors"
	"fmt"

	"github.com/go-telegram/bot"
	telegramModels "github.com/go-telegram/bot/models"
)

type userKeyType string

const userKey userKeyType = "userID"

func ContextUser(ctx context.Context) *models.User {
	value := ctx.Value(userKey)
	user, _ := value.(*models.User)
	return user
}

func ContextWithUser(next bot.HandlerFunc, app *app.App) bot.HandlerFunc {
	return func(ctx context.Context, bot *bot.Bot, update *telegramModels.Update) {
		chatID, err := chatIDFromUpdate(update)
		if err != nil {
			app.Logger.Error("failed extract chatID from update", "error", err)
			return
		}

		user, err := userFromChatID(*chatID, app)
		if err != nil {
			app.Logger.Error("failed get or create user from chat id", "error", err)
			return
		}

		ctx = context.WithValue(ctx, userKey, user)

		next(ctx, bot, update)
	}
}

func chatIDFromUpdate(update *telegramModels.Update) (*int64, error) {
	switch {
	case update.Message != nil:
		return &update.Message.Chat.ID, nil
	case update.CallbackQuery != nil:
		return &update.CallbackQuery.Message.Message.Chat.ID, nil
	default:
		return nil, fmt.Errorf("unable to determine chatID from update")
	}
}

func userFromChatID(chatID int64, app *app.App) (*models.User, error) {
	u, err := app.DB.GetUserByTelegramChatID(chatID)
	if err != nil {
		if errors.Is(err, db.ErrNotExists) {
			u = models.NewUser(chatID)
			u, err = app.DB.CreateUser(u)
			if err != nil {
				return nil, err
			}
			app.Logger.Info("created user", "user", u)
		} else {
			return nil, err
		}
	}
	return u, nil
}
