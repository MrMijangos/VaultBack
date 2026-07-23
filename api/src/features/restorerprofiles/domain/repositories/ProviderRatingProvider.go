package repositories

import "context"

// ProviderRatingProvider calcula la calificacion (1-10) de un proveedor a
// partir de sus reseñas. Mismo patron que RatingProvider en el dominio de
// businesses -- se satisface con PostgreSQLReviewRepository sin depender del
// paquete de reviews.
type ProviderRatingProvider interface {
	GetProviderRating(ctx context.Context, providerID string) (rating *float64, total int, err error)
}
