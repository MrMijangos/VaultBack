package adapters

import (
	"context"
	"fmt"
	"sync"

	"vault-payment/src/features/subscriptions/domain/entities"
)

// InMemorySubscriptionRepository es un stand-in mientras no exista la tabla
// "subscriptions" en Supabase (se deja para el final). Se pierde al
// reiniciar el proceso -- aceptable por ahora, ver nota en payment/README.
type InMemorySubscriptionRepository struct {
	mu          sync.RWMutex
	byID        map[string]*entities.Subscription
	byUserID    map[string]string // userID -> subscriptionID
	byStripeSub map[string]string // stripeSubscriptionID -> subscriptionID
}

func NewInMemorySubscriptionRepository() *InMemorySubscriptionRepository {
	return &InMemorySubscriptionRepository{
		byID:        make(map[string]*entities.Subscription),
		byUserID:    make(map[string]string),
		byStripeSub: make(map[string]string),
	}
}

func (r *InMemorySubscriptionRepository) Create(ctx context.Context, sub *entities.Subscription) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byID[sub.ID] = sub
	r.byUserID[sub.UserID] = sub.ID
	if sub.StripeSubscriptionID != "" {
		r.byStripeSub[sub.StripeSubscriptionID] = sub.ID
	}
	return nil
}

func (r *InMemorySubscriptionRepository) Update(ctx context.Context, sub *entities.Subscription) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byID[sub.ID]; !ok {
		return fmt.Errorf("suscripción %q no existe", sub.ID)
	}
	r.byID[sub.ID] = sub
	if sub.StripeSubscriptionID != "" {
		r.byStripeSub[sub.StripeSubscriptionID] = sub.ID
	}
	return nil
}

func (r *InMemorySubscriptionRepository) GetByUserID(ctx context.Context, userID string) (*entities.Subscription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.byUserID[userID]
	if !ok {
		return nil, nil
	}
	return r.byID[id], nil
}

func (r *InMemorySubscriptionRepository) GetByStripeSubscriptionID(ctx context.Context, stripeSubscriptionID string) (*entities.Subscription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, ok := r.byStripeSub[stripeSubscriptionID]
	if !ok {
		return nil, nil
	}
	return r.byID[id], nil
}
