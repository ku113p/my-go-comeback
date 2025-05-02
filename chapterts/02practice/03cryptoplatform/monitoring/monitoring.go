package monitoring

import (
	"context"
	"crypto/platform/app"
	"crypto/platform/models"
	"crypto/platform/telegram"
)

type Monitoring struct {
	*app.App
	updated <-chan struct{}
}

func NewMonitoring(app *app.App, updated <-chan struct{}) *Monitoring {
	return &Monitoring{app, updated}
}

func (m *Monitoring) Run() error {
	trigger := m.triggerOnUpdate()

	for range trigger {
		if err := m.processUpdate(); err != nil {
			m.Logger.Error("failed processing update", "error", err)
		}
	}

	return nil
}

func (m *Monitoring) triggerOnUpdate() <-chan struct{} {
	trigger := make(chan struct{}, 1)

	go func() {
		for range m.updated {
			select {
			case trigger <- struct{}{}:
			default:
			}
		}
	}()

	return trigger
}

func (m *Monitoring) processUpdate() error {
	users, err := m.DB.ListUsers()
	if err != nil {
		return err
	}

	for _, u := range users {
		go m.notifyUserIfNeed(u, m.App) // TODO many producers and one worker
	}

	return nil
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
	if err := telegram.SendNotification(ctx, n, app); err != nil {
		app.Logger.Error("failed send notification", "error", err)
		return
	}
	app.Logger.Info("sent notification", "notfication", n)
	if err := app.DB.RemoveNotification(*n.ID); err != nil {
		app.Logger.Error("failed delete notification", "error", err)
	}
}
