package request

import "errors"

type UpdateBusinessRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Location    string `json:"location"`
}

func (r UpdateBusinessRequest) Validate() error {
	if r.Name == "" {
		return errors.New("el nombre es obligatorio")
	}
	if !allowedBusinessTypes[r.Type] {
		return errors.New("el tipo de negocio no es valido")
	}
	return nil
}
