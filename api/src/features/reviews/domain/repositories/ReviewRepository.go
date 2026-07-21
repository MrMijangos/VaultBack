package repositories

import (
	"context"
	"errors"

	"vault/src/features/reviews/domain/entities"
)

var ErrReviewNotFound = errors.New("la reseña no existe")

type ReviewRepository interface {
	Create(ctx context.Context, review entities.Review) (entities.Review, error)
	FindByProviderID(ctx context.Context, providerID string) ([]entities.Review, error)
	// GetProviderRating promedia la puntuación (1-10) de las reseñas visibles
	// y ya analizadas de un proveedor. rating es nil si todavía no tiene
	// ninguna reseña analizada.
	GetProviderRating(ctx context.Context, providerID string) (rating *float64, total int, err error)
	FindByID(ctx context.Context, id string) (entities.Review, error)
	Delete(ctx context.Context, id string, userID string) error
	Like(ctx context.Context, reviewID string, userID string) error
	Unlike(ctx context.Context, reviewID string, userID string) error
}
