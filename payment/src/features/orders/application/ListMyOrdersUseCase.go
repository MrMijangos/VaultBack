package application

import (
	"context"

	"vault-payment/src/features/orders/domain/dto/response"
	"vault-payment/src/features/orders/domain/repositories"
)

// ListMyOrdersUseCase alimenta "Mis pedidos" (comprador) -- usa
// ListByBuyerID, que ya existía en OrderRepository sin ningún caso de uso
// que lo llamara.
type ListMyOrdersUseCase struct {
	orderRepo repositories.OrderRepository
}

func NewListMyOrdersUseCase(orderRepo repositories.OrderRepository) *ListMyOrdersUseCase {
	return &ListMyOrdersUseCase{orderRepo: orderRepo}
}

func (uc *ListMyOrdersUseCase) Execute(ctx context.Context, buyerID string) ([]response.OrderResponse, error) {
	orders, err := uc.orderRepo.ListByBuyerID(ctx, buyerID)
	if err != nil {
		return nil, err
	}
	out := make([]response.OrderResponse, 0, len(orders))
	for _, o := range orders {
		out = append(out, response.OrderFromEntity(o))
	}
	return out, nil
}
