package adapters

import (
	"context"
	"fmt"
	"sync"

	"vault-payment/src/features/connect/domain/entities"
)

// InMemoryConnectedAccountRepository -- mismo criterio que el resto de
// payment/: stand-in mientras no exista persistencia real.
type InMemoryConnectedAccountRepository struct {
	mu       sync.RWMutex
	byUserID map[string]*entities.ConnectedAccount
}

func NewInMemoryConnectedAccountRepository() *InMemoryConnectedAccountRepository {
	return &InMemoryConnectedAccountRepository{byUserID: make(map[string]*entities.ConnectedAccount)}
}

func (r *InMemoryConnectedAccountRepository) Create(ctx context.Context, account *entities.ConnectedAccount) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byUserID[account.UserID] = account
	return nil
}

func (r *InMemoryConnectedAccountRepository) Update(ctx context.Context, account *entities.ConnectedAccount) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byUserID[account.UserID]; !ok {
		return fmt.Errorf("el usuario %q no tiene una cuenta conectada", account.UserID)
	}
	r.byUserID[account.UserID] = account
	return nil
}

func (r *InMemoryConnectedAccountRepository) GetByUserID(ctx context.Context, userID string) (*entities.ConnectedAccount, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	account, ok := r.byUserID[userID]
	if !ok {
		return nil, nil
	}
	return account, nil
}
