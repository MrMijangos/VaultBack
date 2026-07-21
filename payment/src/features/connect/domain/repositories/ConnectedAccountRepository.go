package repositories

import (
	"context"

	"vault-payment/src/features/connect/domain/entities"
)

type ConnectedAccountRepository interface {
	Create(ctx context.Context, account *entities.ConnectedAccount) error
	Update(ctx context.Context, account *entities.ConnectedAccount) error
	GetByUserID(ctx context.Context, userID string) (*entities.ConnectedAccount, error)
}
