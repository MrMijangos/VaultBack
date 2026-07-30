package repositories

import "context"

// AssetPriceProvider consulta el precio real del activo en api/ (servicio
// separado -- payment/ no tiene su propia tabla de assets) para no confiar
// en el monto que manda el cliente al crear una orden.
type AssetPriceProvider interface {
	// GetSalePriceCents devuelve el precio de venta vigente del activo en
	// centavos y si sigue publicado (is_for_sale). Un activo que ya no está
	// en venta (vendido, retirado) no debe poder comprarse de nuevo aunque
	// el cliente todavía tenga la pantalla abierta con el precio viejo.
	GetSalePriceCents(ctx context.Context, assetID string) (priceCents int64, isForSale bool, err error)
}
