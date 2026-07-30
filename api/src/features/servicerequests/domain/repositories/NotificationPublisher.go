package repositories

import (
	"context"

	notifentities "vault/src/features/notifications/domain/entities"
)

// NotificationPublisher es el puerto hacia notifications/ -- mismo patrón
// que en comments/chat/posts.domain/repositories.
type NotificationPublisher interface {
	Create(ctx context.Context, notification notifentities.Notification) (notifentities.Notification, error)
}
