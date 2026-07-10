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
	FindByID(ctx context.Context, id string) (entities.Review, error)
	Delete(ctx context.Context, id string, userID string) error
	Like(ctx context.Context, reviewID string, userID string) error
	Unlike(ctx context.Context, reviewID string, userID string) error
}
