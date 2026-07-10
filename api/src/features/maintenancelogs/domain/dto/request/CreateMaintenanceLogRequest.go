package request

import "errors"

var allowedMaintenanceTypes = map[string]bool{
	"mantenimiento": true,
	"restauracion":  true,
}

type CreateMaintenanceLogRequest struct {
	AssetID     string   `json:"asset_id"`
	ProviderID  string   `json:"provider_id"`
	Type        string   `json:"type"`
	Subtype     string   `json:"subtype"`
	Cost        *float64 `json:"cost"`
	PerformedAt string   `json:"performed_at"`
	Notes       string   `json:"notes"`
}

func (r CreateMaintenanceLogRequest) Validate() error {
	if r.AssetID == "" {
		return errors.New("el producto es obligatorio")
	}
	if !allowedMaintenanceTypes[r.Type] {
		return errors.New("el tipo de servicio no es valido")
	}
	return nil
}
