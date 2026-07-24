package request

import "errors"

var allowedBusinessTypes = map[string]bool{
	"restaurador": true,
	"servicio":    true,
}

type CreateBusinessRequest struct {
	Name        string   `json:"name"`
	Types       []string `json:"types"`
	Description string   `json:"description"`
	Location    string   `json:"location"`
	Specialties []string `json:"specialties"`
}

func (r CreateBusinessRequest) Validate() error {
	if r.Name == "" {
		return errors.New("el nombre es obligatorio")
	}
	if len(r.Types) == 0 {
		return errors.New("elige al menos una categoria de negocio")
	}
	for _, t := range r.Types {
		if !allowedBusinessTypes[t] {
			return errors.New("el tipo de negocio no es valido")
		}
	}
	return nil
}
