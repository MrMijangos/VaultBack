package repositories

import "context"

// SellerCommissionProvider es la interfaz angosta hacia "subscriptions"
// (mismo patrón que AdDeactivator/SubscriptionInfoProvider): la comisión de
// Vault depende del plan del vendedor. Si el vendedor no tiene una
// suscripción activa, el adapter decide la tasa por defecto -- orders/ no
// necesita saber esa regla, solo pide "la tasa de este vendedor".
type SellerCommissionProvider interface {
	GetCommissionRate(ctx context.Context, sellerID string) (rate float64, err error)
}
