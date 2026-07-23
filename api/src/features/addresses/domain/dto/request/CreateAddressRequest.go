package request

import "errors"

type CreateAddressRequest struct {
	Label      string `json:"label"`
	Recipient  string `json:"recipient"`
	Phone      string `json:"phone"`
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	References string `json:"references"`
}

func (r CreateAddressRequest) Validate() error {
	if r.Label == "" {
		return errors.New("la etiqueta es obligatoria")
	}
	if r.Recipient == "" {
		return errors.New("el destinatario es obligatorio")
	}
	if r.Phone == "" {
		return errors.New("el telefono es obligatorio")
	}
	if r.Street == "" {
		return errors.New("la calle es obligatoria")
	}
	if r.City == "" {
		return errors.New("la ciudad es obligatoria")
	}
	if r.State == "" {
		return errors.New("el estado es obligatorio")
	}
	if r.PostalCode == "" {
		return errors.New("el codigo postal es obligatorio")
	}
	return nil
}
