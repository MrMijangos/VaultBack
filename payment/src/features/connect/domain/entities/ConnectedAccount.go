package entities

import "time"

// ConnectedAccount es la cuenta de Stripe Connect Express de un vendedor --
// necesaria para poder recibir el Transfer cuando se libera un pago en
// escrow (ver orders/). Se crea la primera vez que el vendedor pide su link
// de onboarding; ChargesEnabled se actualiza consultando a Stripe (el
// vendedor completa el KYC directo en el flujo hospedado de Stripe, no en
// Vault).
type ConnectedAccount struct {
	UserID          string
	StripeAccountID string
	ChargesEnabled  bool
	CreatedAt       time.Time
}
