package repositories

import (
	"context"
	"errors"

	"vault/src/features/notifications/domain/entities"
)

var ErrNotificationNotFound = errors.New("la notificacion no existe")

type NotificationRepository interface {
	Create(ctx context.Context, notification entities.Notification) (entities.Notification, error)
	FindByUserID(ctx context.Context, userID string) ([]entities.Notification, error)
	MarkAsRead(ctx context.Context, id string, userID string) (entities.Notification, error)
	Delete(ctx context.Context, id string, userID string) error
}
