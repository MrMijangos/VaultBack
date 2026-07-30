package repositories

import "context"

// PurchaseVerifier lo satisface un adapter que consulta payment/ (servicio
// separado) para saber si buyerID le compró algo a providerID -- exigir esto
// antes de crear una reseña evita que cualquiera reseñe a cualquier
// vendedor sin haberle comprado nunca.
type PurchaseVerifier interface {
	HasPurchased(ctx context.Context, buyerID string, providerID string) (bool, error)
}
