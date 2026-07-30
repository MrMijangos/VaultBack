package application

import (
	"context"

	"vault-payment/src/features/orders/domain/repositories"
)

type HasPurchasedUseCase struct {
	orderRepo repositories.OrderRepository
}

func NewHasPurchasedUseCase(orderRepo repositories.OrderRepository) *HasPurchasedUseCase {
	return &HasPurchasedUseCase{orderRepo: orderRepo}
}

func (uc *HasPurchasedUseCase) Execute(ctx context.Context, buyerID string, sellerID string) (bool, error) {
	return uc.orderRepo.ExistsCompletedOrder(ctx, buyerID, sellerID)
}
