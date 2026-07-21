package application

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"

	"vault-payment/src/core/stripeclient"
	"vault-payment/src/features/orders/domain/dto/request"
	"vault-payment/src/features/orders/domain/dto/response"
	"vault-payment/src/features/orders/domain/entities"
	"vault-payment/src/features/orders/domain/repositories"
)

// currency -- los planes y precios de Vault están en MXN, ver
// subscriptions/infrastructure/adapters/InMemoryPlanRepository.go.
const currency = "mxn"

type CreateOrderUseCase struct {
	orderRepo          repositories.OrderRepository
	commissionProvider repositories.SellerCommissionProvider
	accountProvider    repositories.SellerAccountProvider
	stripeClient       stripeclient.Client
}

func NewCreateOrderUseCase(
	orderRepo repositories.OrderRepository,
	commissionProvider repositories.SellerCommissionProvider,
	accountProvider repositories.SellerAccountProvider,
	stripeClient stripeclient.Client,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepo:          orderRepo,
		commissionProvider: commissionProvider,
		accountProvider:    accountProvider,
		stripeClient:       stripeClient,
	}
}

func (uc *CreateOrderUseCase) Execute(ctx context.Context, buyerID string, req request.CreateOrderRequest) (*response.OrderResponse, error) {
	if req.SellerID == "" || req.AssetID == "" || req.AmountCents <= 0 || req.BuyerEmail == "" || req.PaymentMethodID == "" {
		return nil, ErrInvalidRequest
	}

	_, ready, err := uc.accountProvider.GetChargesEnabledAccountID(ctx, req.SellerID)
	if err != nil {
		return nil, err
	}
	if !ready {
		return nil, ErrSellerNotOnboarded
	}

	rate, err := uc.commissionProvider.GetCommissionRate(ctx, req.SellerID)
	if err != nil {
		return nil, err
	}

	customerID, attachedPaymentMethodID, err := uc.stripeClient.CreateCustomerWithPaymentMethod(ctx, req.BuyerEmail, req.PaymentMethodID)
	if err != nil {
		return nil, fmt.Errorf("no se pudo registrar el método de pago: %w", err)
	}

	paymentIntentID, err := uc.stripeClient.ChargeBuyer(ctx, customerID, attachedPaymentMethodID, req.AmountCents, currency)
	if err != nil {
		return nil, fmt.Errorf("no se pudo cobrar al comprador: %w", err)
	}

	commissionCents := int64(math.Round(float64(req.AmountCents) * rate))
	sellerAmountCents := req.AmountCents - commissionCents

	order := &entities.Order{
		ID:                    uuid.NewString(),
		BuyerID:               buyerID,
		SellerID:              req.SellerID,
		AssetID:               req.AssetID,
		AmountCents:           req.AmountCents,
		CommissionCents:       commissionCents,
		SellerAmountCents:     sellerAmountCents,
		Currency:              currency,
		Status:                entities.OrderStatusHeld,
		StripeCustomerID:      customerID,
		StripePaymentIntentID: paymentIntentID,
		CreatedAt:             time.Now().UTC(),
	}

	if err := uc.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	out := response.OrderFromEntity(order)
	return &out, nil
}
