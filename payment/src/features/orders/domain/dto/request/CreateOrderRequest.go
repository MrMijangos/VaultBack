package request

// CreateOrderRequest -- payment/ no tiene acceso a la base de datos de
// assets (vive en api/, servicio separado), así que no puede validar el
// precio real del activo por su cuenta todavía; el cliente manda el monto
// ya calculado. Cuando se conecte la persistencia real, esto se valida
// contra el precio guardado del activo antes de cobrar.
type CreateOrderRequest struct {
	SellerID        string `json:"seller_id"`
	AssetID         string `json:"asset_id"`
	AmountCents     int64  `json:"amount_cents"`
	BuyerEmail      string `json:"buyer_email"`
	PaymentMethodID string `json:"payment_method_id"`
}
