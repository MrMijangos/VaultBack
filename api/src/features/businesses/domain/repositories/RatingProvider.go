package repositories

import "context"

// RatingProvider calcula la calificación (1-10) de un proveedor a partir de
// sus reseñas. Vive en el dominio de businesses (no de reviews) para que
// este feature no dependa del paquete de reviews -- en infraestructura se
// satisface con el PostgreSQLReviewRepository, que ya implementa este mismo
// método por su cuenta.
type RatingProvider interface {
	GetProviderRating(ctx context.Context, providerID string) (rating *float64, total int, err error)
}
