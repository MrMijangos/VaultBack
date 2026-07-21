package application

import (
	"context"

	"vault-payment/src/features/subscriptions/domain/dto/response"
	"vault-payment/src/features/subscriptions/domain/repositories"
)

type GetSubscriptionStatusUseCase struct {
	subscriptionRepo repositories.SubscriptionRepository
}

func NewGetSubscriptionStatusUseCase(subscriptionRepo repositories.SubscriptionRepository) *GetSubscriptionStatusUseCase {
	return &GetSubscriptionStatusUseCase{subscriptionRepo: subscriptionRepo}
}

// Execute devuelve nil (sin error) si el usuario nunca se ha suscrito --
// el controller lo traduce a un 200 con subscription: null, no a un 404.
func (uc *GetSubscriptionStatusUseCase) Execute(ctx context.Context, userID string) (*response.SubscriptionResponse, error) {
	sub, err := uc.subscriptionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, nil
	}

	out := response.SubscriptionFromEntity(sub)
	return &out, nil
}
