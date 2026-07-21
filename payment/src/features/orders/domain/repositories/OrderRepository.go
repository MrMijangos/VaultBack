package repositories

import (
	"context"

	"vault-payment/src/features/orders/domain/entities"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entities.Order) error
	Update(ctx context.Context, order *entities.Order) error
	GetByID(ctx context.Context, id string) (*entities.Order, error)
	ListByBuyerID(ctx context.Context, buyerID string) ([]*entities.Order, error)
	ListBySellerID(ctx context.Context, sellerID string) ([]*entities.Order, error)
}
