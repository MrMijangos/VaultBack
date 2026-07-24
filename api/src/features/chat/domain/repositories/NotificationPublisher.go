package repositories

import (
	"context"

	notifentities "vault/src/features/notifications/domain/entities"
)

// NotificationPublisher es el puerto hacia el feature de notifications --
// mismo patrón que RatingProvider en businesses/reviews: en
// infrastructure/dependencies.go se pasa la implementación real de
// notifications/infrastructure/adapters directamente, que ya lo satisface
// por estructura sin que chat dependa de su application/infraestructura.
type NotificationPublisher interface {
	Create(ctx context.Context, notification notifentities.Notification) (notifentities.Notification, error)
}
