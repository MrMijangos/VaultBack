package repositories

import "context"

// SubscriptionInfoProvider es la interfaz angosta hacia "subscriptions"
// (mismo patrón que AdDeactivator, en dirección opuesta): para crear un
// anuncio hay que saber si el usuario tiene una suscripción activa, cuántos
// anuncios le permite su plan y en qué secciones puede publicar.
type SubscriptionInfoProvider interface {
	// GetActiveSubscription devuelve found=false si el usuario no tiene una
	// suscripción activa (nunca se suscribió, o la canceló).
	GetActiveSubscription(ctx context.Context, userID string) (subscriptionID string, maxAds int, targetSections []string, found bool, err error)
}
