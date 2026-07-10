package request

import "errors"

type UpdateMaintenanceLogRequest struct {
	Type        string   `json:"type"`
	Subtype     string   `json:"subtype"`
	Cost        *float64 `json:"cost"`
	PerformedAt string   `json:"performed_at"`
	Notes       string   `json:"notes"`
}

func (r UpdateMaintenanceLogRequest) Validate() error {
	if !allowedMaintenanceTypes[r.Type] {
		return errors.New("el tipo de servicio no es valido")
	}
	return nil
}
