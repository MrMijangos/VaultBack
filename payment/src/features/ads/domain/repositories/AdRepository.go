package repositories

import (
	"context"

	"vault-payment/src/features/ads/domain/entities"
)

type AdRepository interface {
	Create(ctx context.Context, ad *entities.Ad) error
	Update(ctx context.Context, ad *entities.Ad) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*entities.Ad, error)
	ListByUserID(ctx context.Context, userID string) ([]*entities.Ad, error)
	ListActiveBySection(ctx context.Context, section string) ([]*entities.Ad, error)
	// ListBySubscriptionID se usa para desactivar/reactivar todos los
	// anuncios de una suscripción cuando esta se cancela o se renueva.
	ListBySubscriptionID(ctx context.Context, subscriptionID string) ([]*entities.Ad, error)
	CountActiveByUserID(ctx context.Context, userID string) (int, error)
	IncrementImpressions(ctx context.Context, id string) error
	IncrementClicks(ctx context.Context, id string) error
}
