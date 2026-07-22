package request

import "errors"

type UpdateAssetRequest struct {
	Name            string   `json:"name"`
	Category        string   `json:"category"`
	Brand           string   `json:"brand"`
	PurchaseValue   *float64 `json:"purchase_value"`
	Condition       string   `json:"condition"`
	PurchaseDate    string   `json:"purchase_date"`
	StoreOrigin     string   `json:"store_origin"`
	Notes           string   `json:"notes"`
	IsForSale       bool     `json:"is_for_sale"`
	SalePrice       *float64 `json:"sale_price"`
	SaleDescription string   `json:"sale_description"`
	Size            string   `json:"size"`
}

func (r *UpdateAssetRequest) Validate() error {
	if r.Name == "" {
		return errors.New("el nombre es obligatorio")
	}
	if !allowedCategories[r.Category] {
		return errors.New("la categoria no es valida")
	}
	if r.Condition == "" {
		r.Condition = "nuevo"
	}
	if !allowedConditions[r.Condition] {
		return errors.New("la condicion no es valida")
	}
	if r.IsForSale && (r.SalePrice == nil || *r.SalePrice <= 0) {
		return errors.New("el precio de venta es obligatorio para poner el activo en venta")
	}
	return nil
}
