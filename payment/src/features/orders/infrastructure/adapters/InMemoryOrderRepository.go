package adapters

import (
	"context"
	"fmt"
	"sync"

	"vault-payment/src/features/orders/domain/entities"
)

// InMemoryOrderRepository -- mismo criterio que el resto de payment/: stand-in
// mientras no exista persistencia real (ver nota en payment/.env.example).
type InMemoryOrderRepository struct {
	mu   sync.RWMutex
	byID map[string]*entities.Order
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{byID: make(map[string]*entities.Order)}
}

func (r *InMemoryOrderRepository) Create(ctx context.Context, order *entities.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[order.ID] = order
	return nil
}

func (r *InMemoryOrderRepository) Update(ctx context.Context, order *entities.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byID[order.ID]; !ok {
		return fmt.Errorf("orden %q no existe", order.ID)
	}
	r.byID[order.ID] = order
	return nil
}

func (r *InMemoryOrderRepository) GetByID(ctx context.Context, id string) (*entities.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	order, ok := r.byID[id]
	if !ok {
		return nil, nil
	}
	return order, nil
}

func (r *InMemoryOrderRepository) ListByBuyerID(ctx context.Context, buyerID string) ([]*entities.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*entities.Order
	for _, order := range r.byID {
		if order.BuyerID == buyerID {
			out = append(out, order)
		}
	}
	return out, nil
}

func (r *InMemoryOrderRepository) ListBySellerID(ctx context.Context, sellerID string) ([]*entities.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*entities.Order
	for _, order := range r.byID {
		if order.SellerID == sellerID {
			out = append(out, order)
		}
	}
	return out, nil
}
