package adapters

import (
	"context"
	"fmt"
	"sync"

	"vault-payment/src/features/ads/domain/entities"
)

// InMemoryAdRepository -- mismo criterio que InMemorySubscriptionRepository:
// stand-in mientras no exista la tabla "ads" en Supabase.
type InMemoryAdRepository struct {
	mu   sync.RWMutex
	byID map[string]*entities.Ad
}

func NewInMemoryAdRepository() *InMemoryAdRepository {
	return &InMemoryAdRepository{byID: make(map[string]*entities.Ad)}
}

func (r *InMemoryAdRepository) Create(ctx context.Context, ad *entities.Ad) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[ad.ID] = ad
	return nil
}

func (r *InMemoryAdRepository) Update(ctx context.Context, ad *entities.Ad) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byID[ad.ID]; !ok {
		return fmt.Errorf("anuncio %q no existe", ad.ID)
	}
	r.byID[ad.ID] = ad
	return nil
}

func (r *InMemoryAdRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byID[id]; !ok {
		return fmt.Errorf("anuncio %q no existe", id)
	}
	delete(r.byID, id)
	return nil
}

func (r *InMemoryAdRepository) GetByID(ctx context.Context, id string) (*entities.Ad, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ad, ok := r.byID[id]
	if !ok {
		return nil, nil
	}
	return ad, nil
}

func (r *InMemoryAdRepository) ListByUserID(ctx context.Context, userID string) ([]*entities.Ad, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*entities.Ad
	for _, ad := range r.byID {
		if ad.UserID == userID {
			out = append(out, ad)
		}
	}
	return out, nil
}

func (r *InMemoryAdRepository) ListActiveBySection(ctx context.Context, section string) ([]*entities.Ad, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*entities.Ad
	for _, ad := range r.byID {
		if ad.IsActive() && ad.TargetSection == section {
			out = append(out, ad)
		}
	}
	return out, nil
}

func (r *InMemoryAdRepository) ListBySubscriptionID(ctx context.Context, subscriptionID string) ([]*entities.Ad, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*entities.Ad
	for _, ad := range r.byID {
		if ad.SubscriptionID == subscriptionID {
			out = append(out, ad)
		}
	}
	return out, nil
}

func (r *InMemoryAdRepository) IncrementImpressions(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ad, ok := r.byID[id]
	if !ok {
		return fmt.Errorf("anuncio %q no existe", id)
	}
	ad.Impressions++
	return nil
}

func (r *InMemoryAdRepository) IncrementClicks(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ad, ok := r.byID[id]
	if !ok {
		return fmt.Errorf("anuncio %q no existe", id)
	}
	ad.Clicks++
	return nil
}

func (r *InMemoryAdRepository) CountActiveByUserID(ctx context.Context, userID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, ad := range r.byID {
		if ad.UserID == userID && ad.IsActive() {
			count++
		}
	}
	return count, nil
}
