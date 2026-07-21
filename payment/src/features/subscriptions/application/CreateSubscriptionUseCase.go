package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"vault-payment/src/core/eventbus"
	"vault-payment/src/core/stripeclient"
	"vault-payment/src/features/subscriptions/domain/dto/request"
	"vault-payment/src/features/subscriptions/domain/dto/response"
	"vault-payment/src/features/subscriptions/domain/entities"
	"vault-payment/src/features/subscriptions/domain/repositories"
)

type CreateSubscriptionUseCase struct {
	planRepo         repositories.PlanRepository
	subscriptionRepo repositories.SubscriptionRepository
	stripeClient     stripeclient.Client
	publisher        eventbus.Publisher
}

func NewCreateSubscriptionUseCase(
	planRepo repositories.PlanRepository,
	subscriptionRepo repositories.SubscriptionRepository,
	stripeClient stripeclient.Client,
	publisher eventbus.Publisher,
) *CreateSubscriptionUseCase {
	return &CreateSubscriptionUseCase{
		planRepo:         planRepo,
		subscriptionRepo: subscriptionRepo,
		stripeClient:     stripeClient,
		publisher:        publisher,
	}
}

func (uc *CreateSubscriptionUseCase) Execute(ctx context.Context, userID, role string, req request.CreateSubscriptionRequest) (*response.SubscriptionResponse, error) {
	if !isRoleAllowed(role) {
		return nil, ErrRoleNotAllowed
	}
	if req.PlanID == "" || req.Email == "" || req.PaymentMethodID == "" {
		return nil, ErrInvalidRequest
	}

	existing, err := uc.subscriptionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existing != nil && existing.IsActive() {
		return nil, ErrAlreadySubscribed
	}

	plan, err := uc.planRepo.GetByID(ctx, req.PlanID)
	if err != nil {
		return nil, err
	}

	customerID, _, err := uc.stripeClient.CreateCustomerWithPaymentMethod(ctx, req.Email, req.PaymentMethodID)
	if err != nil {
		return nil, fmt.Errorf("no se pudo registrar el método de pago: %w", err)
	}

	stripeSubID, periodStart, periodEnd, err := uc.stripeClient.CreateSubscription(ctx, customerID, plan.StripePriceID)
	if err != nil {
		return nil, fmt.Errorf("no se pudo crear la suscripción en Stripe: %w", err)
	}

	sub := &entities.Subscription{
		ID:                   uuid.NewString(),
		UserID:               userID,
		PlanID:               plan.ID,
		Status:               entities.SubscriptionStatusActive,
		StripeCustomerID:     customerID,
		StripeSubscriptionID: stripeSubID,
		CurrentPeriodStart:   periodStart,
		CurrentPeriodEnd:     periodEnd,
		CreatedAt:            time.Now().UTC(),
	}

	if err := uc.subscriptionRepo.Create(ctx, sub); err != nil {
		return nil, err
	}

	_ = uc.publisher.PublishSubscriptionEvent(ctx, eventbus.SubscriptionEventPayload{
		EventType:      eventbus.EventSubscriptionActivated,
		UserID:         userID,
		SubscriptionID: sub.ID,
		PlanName:       plan.Name,
	})

	out := response.SubscriptionFromEntity(sub)
	return &out, nil
}
