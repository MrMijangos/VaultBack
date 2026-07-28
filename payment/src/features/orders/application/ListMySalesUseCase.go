package application

import (
	"context"

	"vault-payment/src/features/orders/domain/dto/response"
	"vault-payment/src/features/orders/domain/repositories"
)

// ListMySalesUseCase alimenta "Mis ventas" (vendedor) -- usa
// ListBySellerID, que ya existía en OrderRepository sin ningún caso de uso
// que lo llamara.
type ListMySalesUseCase struct {
	orderRepo repositories.OrderRepository
}

func NewListMySalesUseCase(orderRepo repositories.OrderRepository) *ListMySalesUseCase {
	return &ListMySalesUseCase{orderRepo: orderRepo}
}

func (uc *ListMySalesUseCase) Execute(ctx context.Context, sellerID string) ([]response.OrderResponse, error) {
	orders, err := uc.orderRepo.ListBySellerID(ctx, sellerID)
	if err != nil {
		return nil, err
	}
	out := make([]response.OrderResponse, 0, len(orders))
	for _, o := range orders {
		out = append(out, response.OrderFromEntity(o))
	}
	return out, nil
}
