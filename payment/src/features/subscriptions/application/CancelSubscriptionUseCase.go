package application

import (
	"context"
	"fmt"
	"time"

	"vault-payment/src/core/eventbus"
	"vault-payment/src/core/stripeclient"
	"vault-payment/src/features/subscriptions/domain/entities"
	"vault-payment/src/features/subscriptions/domain/repositories"
)

type CancelSubscriptionUseCase struct {
	subscriptionRepo repositories.SubscriptionRepository
	stripeClient     stripeclient.Client
	adDeactivator    repositories.AdDeactivator
	publisher        eventbus.Publisher
}

func NewCancelSubscriptionUseCase(
	subscriptionRepo repositories.SubscriptionRepository,
	stripeClient stripeclient.Client,
	adDeactivator repositories.AdDeactivator,
	publisher eventbus.Publisher,
) *CancelSubscriptionUseCase {
	return &CancelSubscriptionUseCase{
		subscriptionRepo: subscriptionRepo,
		stripeClient:     stripeClient,
		adDeactivator:    adDeactivator,
		publisher:        publisher,
	}
}

func (uc *CancelSubscriptionUseCase) Execute(ctx context.Context, userID string) error {
	sub, err := uc.subscriptionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if sub == nil || !sub.IsActive() {
		return ErrNotSubscribed
	}

	if err := uc.stripeClient.CancelSubscription(ctx, sub.StripeSubscriptionID); err != nil {
		return fmt.Errorf("no se pudo cancelar la suscripción en Stripe: %w", err)
	}

	now := time.Now().UTC()
	sub.Status = entities.SubscriptionStatusCanceled
	sub.CanceledAt = &now
	if err := uc.subscriptionRepo.Update(ctx, sub); err != nil {
		return err
	}

	if err := uc.adDeactivator.DeactivateBySubscriptionID(ctx, sub.ID); err != nil {
		return fmt.Errorf("la suscripción se canceló pero no se pudieron desactivar sus anuncios: %w", err)
	}

	_ = uc.publisher.PublishSubscriptionEvent(ctx, eventbus.SubscriptionEventPayload{
		EventType:      eventbus.EventSubscriptionCanceled,
		UserID:         userID,
		SubscriptionID: sub.ID,
	})

	return nil
}
