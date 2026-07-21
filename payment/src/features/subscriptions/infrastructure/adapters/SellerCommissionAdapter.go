package adapters

import (
	"context"

	"vault-payment/src/features/subscriptions/domain/entities"
	"vault-payment/src/features/subscriptions/domain/repositories"
)

// SellerCommissionAdapter implementa orders/domain/repositories.SellerCommissionProvider
// (tipado estructural, sin import cruzado entre features).
type SellerCommissionAdapter struct {
	subscriptionRepo repositories.SubscriptionRepository
	planRepo         repositories.PlanRepository
}

func NewSellerCommissionAdapter(subscriptionRepo repositories.SubscriptionRepository, planRepo repositories.PlanRepository) *SellerCommissionAdapter {
	return &SellerCommissionAdapter{subscriptionRepo: subscriptionRepo, planRepo: planRepo}
}

func (a *SellerCommissionAdapter) GetCommissionRate(ctx context.Context, sellerID string) (float64, error) {
	sub, err := a.subscriptionRepo.GetByUserID(ctx, sellerID)
	if err != nil {
		return 0, err
	}
	if sub == nil || !sub.IsActive() {
		return entities.DefaultCommissionRate, nil
	}

	plan, err := a.planRepo.GetByID(ctx, sub.PlanID)
	if err != nil {
		return 0, err
	}
	return plan.CommissionRate, nil
}
