package request

import "errors"

var allowedServiceRequestTypes = map[string]bool{
	"servicio":   true,
	"reparacion": true,
}

type CreateServiceRequestRequest struct {
	AssetID    string `json:"asset_id"`
	BusinessID string `json:"business_id"`
	Type       string `json:"type"`
}

func (r CreateServiceRequestRequest) Validate() error {
	if r.AssetID == "" {
		return errors.New("el activo es obligatorio")
	}
	if r.BusinessID == "" {
		return errors.New("el negocio es obligatorio")
	}
	if !allowedServiceRequestTypes[r.Type] {
		return errors.New("el tipo de servicio no es valido")
	}
	return nil
}
