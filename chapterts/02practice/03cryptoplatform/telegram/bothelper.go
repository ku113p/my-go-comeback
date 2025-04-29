package telegram

import (
	"context"
	"crypto/platform/app"
	"crypto/platform/db"
	"crypto/platform/models"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	telegramModels "github.com/go-telegram/bot/models"
)

type BotHelper struct {
	app  *app.App
	mode mode
}

func (h *BotHelper) Run() error {
	opts := []bot.Option{
		bot.WithMiddlewares(h.withUser),
		bot.WithDefaultHandler(h.defaultHandler()),
		bot.WithMessageTextHandler("help", bot.MatchTypeCommand, helpCommand),
		bot.WithMessageTextHandler("add", bot.MatchTypeCommand, h.addCommand()),
	}

	return h.mode.runBot(opts...)
}

func (h *BotHelper) defaultHandler() func(context.Context, *bot.Bot, *telegramModels.Update) {
	return func(ctx context.Context, b *bot.Bot, update *telegramModels.Update) {
		if update.Message != nil {
			u := user(ctx)
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("%#v", u),
			})
		}
	}
}

func helpCommand(ctx context.Context, b *bot.Bot, update *telegramModels.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "This bot help to monitor crypto prices",
	})
}

func (h *BotHelper) withUser(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *telegramModels.Update) {
		withUser(ctx, b, update, next, h.app)
	}
}

type userKeyType string

const userKey userKeyType = "user"

func withUser(ctx context.Context, b *bot.Bot, update *telegramModels.Update, next bot.HandlerFunc, a *app.App) {
	msg := update.Message
	if msg != nil {
		if telegramUser := msg.From; telegramUser != nil {
			u, err := a.DB.GetUserByTelegramChatID(telegramUser.ID)
			if err != nil {
				if errors.Is(err, db.ErrNotExists) {
					u = models.NewUser(telegramUser.ID)
					u, err = a.DB.CreateUser(u)
					if err != nil {
						a.Logger.Error("failed create user", "error", err)
						return
					}
					a.Logger.Info("created user", "user", u)
				} else {
					a.Logger.Error("failed get user from db", "error", err)
					return
				}
			}
			ctx = context.WithValue(ctx, userKey, u)
		}
	}
	next(ctx, b, update)
}

func user(ctx context.Context) *models.User {
	value := ctx.Value(userKey)
	u, _ := value.(*models.User)
	return u
}

func (h *BotHelper) addCommand() func(context.Context, *bot.Bot, *telegramModels.Update) {
	return func(ctx context.Context, b *bot.Bot, update *telegramModels.Update) {
		addCommand(ctx, b, update, h.app)
	}
}

type addSign string

const (
	moreSign addSign = ">"
	lessSign addSign = "<"
)

func getSign(s string) (*addSign, error) {
	sign := addSign(s)
	switch sign {
	case moreSign, lessSign:
		return &sign, nil
	default:
		return nil, errors.New("invalid sign")
	}
}

func (s *addSign) checkFunction(symbol string, amount float64) func(p *models.TokenPrice) bool {
	return func(p *models.TokenPrice) bool {
		if p.Symbol != symbol {
			return false
		}

		switch *s {
		case moreSign:
			return p.Price > amount
		case lessSign:
			return p.Price < amount
		}

		return false
	}
}

func (s *addSign) String() string {
	switch *s {
	case moreSign:
		return ">"
	case lessSign:
		return "<"
	}

	return "?"
}

func addCommand(ctx context.Context, b *bot.Bot, update *telegramModels.Update, a *app.App) {
	u := user(ctx)
	sendError := func(msg string) { sendErrorMessage(ctx, update.Message.Chat.ID, msg, b) }

	n, err := newNotificationFromString(update.Message.Text)
	if err != nil {
		sendError(fmt.Sprintf("%s", err))
		return
	}

	token, err := a.DB.GetPrice(n.Symbol)
	if err != nil {
		sendError("")
		a.Logger.Error("failed get price", "error", err)
		return
	}

	if n.Check(token) {
		sendError("price already reached target amount")
		return
	}

	n.UserID = u.ID
	n, err = a.DB.CreateNotification(n)
	if err != nil {
		sendError("")
		a.Logger.Error("failed create notification", "error", err)
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Notification #{%s} created.", *n.ID),
	})
}

func newNotificationFromString(s string) (*models.Notification, error) {
	s = strings.Replace(s, "/add ", "", 1)
	s = strings.Trim(s, " ")

	words := strings.SplitN(s, " ", 3)
	if len(words) != 3 {
		return nil, fmt.Errorf("invalid format")
	}

	symbol, signString, amountString := words[0], words[1], words[2]
	symbol = strings.ToUpper(symbol)

	sign, err := getSign(signString)
	if err != nil {
		return nil, fmt.Errorf("invalid sign")
	}

	amountString = strings.ReplaceAll(amountString, ",", ".")
	amount, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid amount")
	}

	checkFuncion := sign.checkFunction(symbol, amount)

	msg := fmt.Sprintf("price %v %v %v", symbol, sign.String(), amount)
	n := models.NewNotification(symbol, checkFuncion, nil, &msg)

	return n, nil
}

func sendErrorMessage(ctx context.Context, chatID int64, msg string, b *bot.Bot) error {
	if msg == "" {
		msg = "Unexpected error"
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   msg,
	})

	return err
}
