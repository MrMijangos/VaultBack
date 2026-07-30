package request

// CreateOrderRequest -- AmountCents ya NO se usa para cobrar (ver
// CreateOrderUseCase.Execute): el monto real se consulta en api/ vía
// AssetPriceProvider, para que el cliente no pueda decidir cuánto paga por
// un activo. Se deja el campo para no romper clientes que lo sigan mandando
// -- simplemente se ignora.
type CreateOrderRequest struct {
	SellerID        string `json:"seller_id"`
	AssetID         string `json:"asset_id"`
	AmountCents     int64  `json:"amount_cents"`
	BuyerEmail      string `json:"buyer_email"`
	PaymentMethodID string `json:"payment_method_id"`
}
