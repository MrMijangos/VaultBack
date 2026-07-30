package repositories

import "context"

// PurchaseVerifier es un port hacia el servicio de pagos -- mismo patrón que
// reviews/domain/repositories.PurchaseVerifier, cada feature declara su
// propia interfaz mínima aunque el adaptador concreto (HTTPPurchaseVerifier)
// sea compartido.
type PurchaseVerifier interface {
	HasPurchased(ctx context.Context, buyerID, sellerID string) (bool, error)
}
