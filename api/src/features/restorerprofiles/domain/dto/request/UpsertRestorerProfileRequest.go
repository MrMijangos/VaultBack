package request

import "errors"

type RestorerServiceRequest struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type UpsertRestorerProfileRequest struct {
	Bio         string                   `json:"bio"`
	Specialties []string                 `json:"specialties"`
	Services    []RestorerServiceRequest `json:"services"`
}

func (r UpsertRestorerProfileRequest) Validate() error {
	for _, s := range r.Services {
		if s.Title == "" {
			return errors.New("el titulo del servicio es obligatorio")
		}
		if s.Price <= 0 {
			return errors.New("el precio del servicio debe ser mayor a cero")
		}
	}
	return nil
}
