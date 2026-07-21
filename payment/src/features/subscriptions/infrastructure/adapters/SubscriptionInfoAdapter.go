package adapters

import (
	"context"

	"vault-payment/src/features/subscriptions/domain/repositories"
)

// SubscriptionInfoAdapter implementa ads/domain/repositories.SubscriptionInfoProvider
// (tipado estructural, sin import cruzado entre features).
type SubscriptionInfoAdapter struct {
	subscriptionRepo repositories.SubscriptionRepository
	planRepo         repositories.PlanRepository
}

func NewSubscriptionInfoAdapter(subscriptionRepo repositories.SubscriptionRepository, planRepo repositories.PlanRepository) *SubscriptionInfoAdapter {
	return &SubscriptionInfoAdapter{subscriptionRepo: subscriptionRepo, planRepo: planRepo}
}

func (a *SubscriptionInfoAdapter) GetActiveSubscription(ctx context.Context, userID string) (string, int, []string, bool, error) {
	sub, err := a.subscriptionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return "", 0, nil, false, err
	}
	if sub == nil || !sub.IsActive() {
		return "", 0, nil, false, nil
	}

	plan, err := a.planRepo.GetByID(ctx, sub.PlanID)
	if err != nil {
		return "", 0, nil, false, err
	}

	return sub.ID, plan.MaxAds, plan.TargetSections, true, nil
}
