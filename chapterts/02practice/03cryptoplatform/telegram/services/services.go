package services

import "crypto/platform/app"

type Services struct {
	Notification *NotificationService
}

func NewServices(app *app.App) *Services {
	n := newNotificationService(app)

	return &Services{n}
}
