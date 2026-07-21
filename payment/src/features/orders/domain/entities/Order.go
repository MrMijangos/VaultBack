package entities

import "time"

const (
	// OrderStatusHeld: el comprador ya pagó, el dinero está retenido en el
	// balance de Stripe de Vault, esperando a que el comprador confirme
	// que recibió el producto.
	OrderStatusHeld = "retenido"
	// OrderStatusReleased: el comprador confirmó recibido -- el dinero (menos
	// la comisión de Vault) ya se transfirió a la cuenta Stripe del vendedor.
	OrderStatusReleased = "liberado"
)

// Order es una compra-venta con el pago retenido en escrow -- distinto de
// Subscription (suscripción del vendedor a Vault): aquí el dinero es del
// comprador y eventualmente le pertenece al vendedor, Vault solo lo retiene
// y se queda con una comisión al liberarlo.
type Order struct {
	ID                    string
	BuyerID               string
	SellerID              string
	AssetID               string
	AmountCents           int64
	CommissionCents       int64
	SellerAmountCents     int64
	Currency              string
	Status                string
	StripeCustomerID      string
	StripePaymentIntentID string
	StripeTransferID      string
	CreatedAt             time.Time
	ConfirmedAt           *time.Time
}
