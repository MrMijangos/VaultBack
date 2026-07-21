package application

import (
	"context"
	"fmt"
	"time"

	"vault-payment/src/core/eventbus"
	"vault-payment/src/core/stripeclient"
	"vault-payment/src/features/orders/domain/dto/response"
	"vault-payment/src/features/orders/domain/entities"
	"vault-payment/src/features/orders/domain/repositories"
)

// ConfirmOrderUseCase libera el escrow: transfiere el dinero (menos la
// comisión de Vault) al vendedor. api/ todavía no tiene el consumidor que
// convierte EventOrderConfirmed en una transferencia de propiedad
// blockchain -- se agrega junto con la persistencia real (ver spec de
// suscripciones sección 2).
type ConfirmOrderUseCase struct {
	orderRepo       repositories.OrderRepository
	accountProvider repositories.SellerAccountProvider
	stripeClient    stripeclient.Client
	publisher       eventbus.Publisher
}

func NewConfirmOrderUseCase(
	orderRepo repositories.OrderRepository,
	accountProvider repositories.SellerAccountProvider,
	stripeClient stripeclient.Client,
	publisher eventbus.Publisher,
) *ConfirmOrderUseCase {
	return &ConfirmOrderUseCase{
		orderRepo:       orderRepo,
		accountProvider: accountProvider,
		stripeClient:    stripeClient,
		publisher:       publisher,
	}
}

func (uc *ConfirmOrderUseCase) Execute(ctx context.Context, buyerID, orderID string) (*response.OrderResponse, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}
	if order.BuyerID != buyerID {
		return nil, ErrNotBuyer
	}
	if order.Status != entities.OrderStatusHeld {
		return nil, ErrNotHeld
	}

	accountID, ready, err := uc.accountProvider.GetChargesEnabledAccountID(ctx, order.SellerID)
	if err != nil {
		return nil, err
	}
	if !ready {
		return nil, ErrSellerNotOnboarded
	}

	transferID, err := uc.stripeClient.ReleaseToSeller(ctx, accountID, order.SellerAmountCents, order.Currency, order.ID)
	if err != nil {
		return nil, fmt.Errorf("no se pudo transferir el pago al vendedor: %w", err)
	}

	now := time.Now().UTC()
	order.Status = entities.OrderStatusReleased
	order.StripeTransferID = transferID
	order.ConfirmedAt = &now

	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	_ = uc.publisher.PublishOrderEvent(ctx, eventbus.OrderEventPayload{
		EventType: eventbus.EventOrderConfirmed,
		OrderID:   order.ID,
		BuyerID:   order.BuyerID,
		SellerID:  order.SellerID,
		AssetID:   order.AssetID,
	})

	out := response.OrderFromEntity(order)
	return &out, nil
}
