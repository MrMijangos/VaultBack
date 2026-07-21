package repositories

import "context"

// AdDeactivator es una interfaz angosta hacia la feature "ads" (mismo
// patrón que RatingProvider en api/src/features/businesses): al cancelar
// una suscripción (por API o por webhook de Stripe) hay que desactivar sus
// anuncios sin que "subscriptions" importe el paquete de "ads" directamente.
type AdDeactivator interface {
	DeactivateBySubscriptionID(ctx context.Context, subscriptionID string) error
}
