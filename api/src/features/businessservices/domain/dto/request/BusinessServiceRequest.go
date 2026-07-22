package request

import "errors"

type BusinessServiceRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func (r BusinessServiceRequest) Validate() error {
	if r.Title == "" {
		return errors.New("el titulo es obligatorio")
	}
	if r.Price <= 0 {
		return errors.New("el precio debe ser mayor a cero")
	}
	return nil
}
