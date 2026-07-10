package request

import "errors"

var allowedCategories = map[string]bool{
	"sneakers":       true,
	"gorras":         true,
	"relojes":        true,
	"lentes":         true,
	"carteras":       true,
	"bolsos":         true,
	"pulsos":         true,
	"bisuteria":      true,
	"coleccionables": true,
	"otros":          true,
}

var allowedConditions = map[string]bool{
	"nuevo":     true,
	"seminuevo": true,
	"usado":     true,
}

type CreateAssetRequest struct {
	Name          string   `json:"name"`
	Category      string   `json:"category"`
	Brand         string   `json:"brand"`
	PurchaseValue *float64 `json:"purchase_value"`
	Condition     string   `json:"condition"`
	PurchaseDate  string   `json:"purchase_date"`
	StoreOrigin   string   `json:"store_origin"`
	Notes         string   `json:"notes"`
}

func (r *CreateAssetRequest) Validate() error {
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
	return nil
}
