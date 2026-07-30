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
	// ExistsCompletedOrder indica si buyerID tiene al menos un pedido en
	// estado "liberado" con sellerID -- lo usa reviews/ (api/) para exigir
	// una compra completada antes de dejar reseñar a un vendedor.
	ExistsCompletedOrder(ctx context.Context, buyerID string, sellerID string) (bool, error)
}
