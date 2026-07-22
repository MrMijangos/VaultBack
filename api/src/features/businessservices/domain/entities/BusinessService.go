package entities

import "time"

// BusinessService es un item del catálogo de servicios que ofrece un negocio
// tipo "restaurador" o "servicio" (ej. "Limpieza de sneakers", "Ajuste de
// correa"), con su precio. Vive bajo un business_id.
type BusinessService struct {
	ID          string
	BusinessID  string
	Title       string
	Description string
	Price       float64
	CreatedAt   time.Time
}
