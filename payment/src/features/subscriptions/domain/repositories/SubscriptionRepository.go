package repositories

import (
	"context"

	"vault-payment/src/features/subscriptions/domain/entities"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *entities.Subscription) error
	Update(ctx context.Context, sub *entities.Subscription) error
	GetByUserID(ctx context.Context, userID string) (*entities.Subscription, error)
	GetByStripeSubscriptionID(ctx context.Context, stripeSubscriptionID string) (*entities.Subscription, error)
}
