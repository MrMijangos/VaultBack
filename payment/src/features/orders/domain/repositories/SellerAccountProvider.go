package repositories

import "context"

// SellerAccountProvider es la interfaz angosta hacia "connect": para poder
// vender, el vendedor necesita una cuenta de Stripe Connect con
// charges_enabled=true -- si no, Stripe no tiene a dónde transferirle el
// pago liberado.
type SellerAccountProvider interface {
	GetChargesEnabledAccountID(ctx context.Context, sellerID string) (stripeAccountID string, ready bool, err error)
}
