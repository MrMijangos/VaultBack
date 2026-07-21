package application

import (
	"context"

	"vault-payment/src/features/orders/domain/dto/response"
	"vault-payment/src/features/orders/domain/repositories"
)

type GetOrderUseCase struct {
	orderRepo repositories.OrderRepository
}

func NewGetOrderUseCase(orderRepo repositories.OrderRepository) *GetOrderUseCase {
	return &GetOrderUseCase{orderRepo: orderRepo}
}

// Execute solo deja ver la orden al comprador o al vendedor involucrados.
func (uc *GetOrderUseCase) Execute(ctx context.Context, userID, orderID string) (*response.OrderResponse, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}
	if order.BuyerID != userID && order.SellerID != userID {
		return nil, ErrNotBuyer
	}

	out := response.OrderFromEntity(order)
	return &out, nil
}
