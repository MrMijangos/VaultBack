package application

import "errors"

var (
	ErrInvalidRequest     = errors.New("seller_id, asset_id, amount_cents, buyer_email y payment_method_id son obligatorios")
	ErrSellerNotOnboarded = errors.New("el vendedor todavía no completó el onboarding de pagos")
	ErrOrderNotFound      = errors.New("orden no encontrada")
	ErrNotBuyer           = errors.New("esta orden no te pertenece")
	ErrNotHeld            = errors.New("esta orden ya fue liberada o no está retenida")
)
