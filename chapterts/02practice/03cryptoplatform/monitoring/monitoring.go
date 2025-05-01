package monitoring

import (
	"context"
	"crypto/platform/app"
	"crypto/platform/models"
	"crypto/platform/telegram"
)

type Monitoring struct {
	*app.App
	updated <-chan any
}

func NewMonitoring(app *app.App, updated <-chan any) *Monitoring {
	return &Monitoring{app, updated}
}

func (m *Monitoring) Run() error {
	for {
		<-m.updated

		users, err := m.DB.ListUsers()
		if err != nil {
			m.Logger.Error("failed get users", "error", err)
			continue
		}

		for _, u := range users {
			go m.notifyUserIfNeed(u, m.App)
		}
	}
}

func (m *Monitoring) notifyUserIfNeed(u *models.User, app *app.App) error {
	ns, err := app.DB.ListNotificationsByUserID(*u.ID)
	if err != nil {
		return err
	}

	for _, n := range ns {
		token, err := app.DB.GetPrice(n.Symbol)
		if err != nil {
			app.Logger.Error("failed get price", "error", err)
		}
		if n.Check(token) {
			go sendNotification(context.TODO(), n, app)
		}
	}

	return nil
}

func sendNotification(ctx context.Context, n *models.Notification, app *app.App) {
	telegram.SendNotification(ctx, n, app)
	app.Logger.Info("sent notification", "notfication", n)
	if err := app.DB.RemoveNotification(*n.ID); err != nil {
		app.Logger.Error("failed delete notification", "error", err)
	}
}
