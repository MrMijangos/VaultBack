package application

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/stripe/stripe-go/v86"

	"vault-payment/src/core/eventbus"
	"vault-payment/src/core/stripeclient"
	"vault-payment/src/features/subscriptions/domain/entities"
	"vault-payment/src/features/subscriptions/domain/repositories"
)

// HandleStripeWebhookUseCase reacciona a eventos que pueden originarse fuera
// de esta API (ej. el usuario cancela desde el dashboard de Stripe, o un
// cobro recurrente falla) -- por eso la cancelación desde CancelSubscriptionUseCase
// y este webhook comparten la misma lógica de desactivar anuncios.
type HandleStripeWebhookUseCase struct {
	subscriptionRepo repositories.SubscriptionRepository
	stripeClient     stripeclient.Client
	adDeactivator    repositories.AdDeactivator
	publisher        eventbus.Publisher
	webhookSecret    string
}

func NewHandleStripeWebhookUseCase(
	subscriptionRepo repositories.SubscriptionRepository,
	stripeClient stripeclient.Client,
	adDeactivator repositories.AdDeactivator,
	publisher eventbus.Publisher,
	webhookSecret string,
) *HandleStripeWebhookUseCase {
	return &HandleStripeWebhookUseCase{
		subscriptionRepo: subscriptionRepo,
		stripeClient:     stripeClient,
		adDeactivator:    adDeactivator,
		publisher:        publisher,
		webhookSecret:    webhookSecret,
	}
}

func (uc *HandleStripeWebhookUseCase) Execute(ctx context.Context, payload []byte, sigHeader string) error {
	event, err := uc.stripeClient.ConstructWebhookEvent(payload, sigHeader, uc.webhookSecret)
	if err != nil {
		return fmt.Errorf("firma de webhook inválida: %w", err)
	}

	switch event.Type {
	case stripe.EventTypeCustomerSubscriptionDeleted:
		return uc.handleSubscriptionDeleted(ctx, event)
	case stripe.EventTypeCustomerSubscriptionUpdated:
		return uc.handleSubscriptionUpdated(ctx, event)
	case stripe.EventTypeInvoicePaymentFailed:
		return uc.handleInvoicePaymentFailed(ctx, event)
	default:
		return nil
	}
}

func (uc *HandleStripeWebhookUseCase) handleSubscriptionDeleted(ctx context.Context, event stripe.Event) error {
	var stripeSub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &stripeSub); err != nil {
		return fmt.Errorf("no se pudo leer el objeto subscription del webhook: %w", err)
	}

	sub, err := uc.subscriptionRepo.GetByStripeSubscriptionID(ctx, stripeSub.ID)
	if err != nil {
		return err
	}
	if sub == nil {
		return nil
	}

	sub.Status = entities.SubscriptionStatusCanceled
	if err := uc.subscriptionRepo.Update(ctx, sub); err != nil {
		return err
	}

	if err := uc.adDeactivator.DeactivateBySubscriptionID(ctx, sub.ID); err != nil {
		return err
	}

	return uc.publisher.PublishSubscriptionEvent(ctx, eventbus.SubscriptionEventPayload{
		EventType:      eventbus.EventSubscriptionCanceled,
		UserID:         sub.UserID,
		SubscriptionID: sub.ID,
	})
}

func (uc *HandleStripeWebhookUseCase) handleSubscriptionUpdated(ctx context.Context, event stripe.Event) error {
	var stripeSub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &stripeSub); err != nil {
		return fmt.Errorf("no se pudo leer el objeto subscription del webhook: %w", err)
	}

	sub, err := uc.subscriptionRepo.GetByStripeSubscriptionID(ctx, stripeSub.ID)
	if err != nil {
		return err
	}
	if sub == nil {
		return nil
	}

	if stripeSub.Items != nil && len(stripeSub.Items.Data) > 0 {
		item := stripeSub.Items.Data[0]
		sub.CurrentPeriodStart = time.Unix(item.CurrentPeriodStart, 0).UTC()
		sub.CurrentPeriodEnd = time.Unix(item.CurrentPeriodEnd, 0).UTC()
	}
	if stripeSub.Status == stripe.SubscriptionStatusActive {
		sub.Status = entities.SubscriptionStatusActive
	}

	if err := uc.subscriptionRepo.Update(ctx, sub); err != nil {
		return err
	}

	return uc.publisher.PublishSubscriptionEvent(ctx, eventbus.SubscriptionEventPayload{
		EventType:      eventbus.EventSubscriptionRenewed,
		UserID:         sub.UserID,
		SubscriptionID: sub.ID,
	})
}

func (uc *HandleStripeWebhookUseCase) handleInvoicePaymentFailed(ctx context.Context, event stripe.Event) error {
	var invoice stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		return fmt.Errorf("no se pudo leer el objeto invoice del webhook: %w", err)
	}
	if invoice.Parent == nil || invoice.Parent.SubscriptionDetails == nil || invoice.Parent.SubscriptionDetails.Subscription == nil {
		return nil
	}

	sub, err := uc.subscriptionRepo.GetByStripeSubscriptionID(ctx, invoice.Parent.SubscriptionDetails.Subscription.ID)
	if err != nil {
		return err
	}
	if sub == nil {
		return nil
	}

	sub.Status = entities.SubscriptionStatusFailed
	if err := uc.subscriptionRepo.Update(ctx, sub); err != nil {
		return err
	}

	return uc.publisher.PublishSubscriptionEvent(ctx, eventbus.SubscriptionEventPayload{
		EventType:      eventbus.EventSubscriptionFailed,
		UserID:         sub.UserID,
		SubscriptionID: sub.ID,
	})
}
