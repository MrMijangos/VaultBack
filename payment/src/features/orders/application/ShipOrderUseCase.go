package application

import (
	"context"

	"vault-payment/src/core/eventbus"
	"vault-payment/src/features/orders/domain/dto/response"
	"vault-payment/src/features/orders/domain/entities"
	"vault-payment/src/features/orders/domain/repositories"
)

// ShipOrderUseCase es el paso de logística del vendedor: marca la orden
// como enviada, requisito para que el comprador pueda confirmar recibido
// (ver ConfirmOrderUseCase, que ahora exige OrderStatusShipped en vez de
// Held).
type ShipOrderUseCase struct {
	orderRepo repositories.OrderRepository
	publisher eventbus.Publisher
}

func NewShipOrderUseCase(orderRepo repositories.OrderRepository, publisher eventbus.Publisher) *ShipOrderUseCase {
	return &ShipOrderUseCase{orderRepo: orderRepo, publisher: publisher}
}

func (uc *ShipOrderUseCase) Execute(ctx context.Context, sellerID, orderID string) (*response.OrderResponse, error) {
	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound
	}
	if order.SellerID != sellerID {
		return nil, ErrNotSeller
	}
	if order.Status != entities.OrderStatusHeld {
		return nil, ErrNotHeld
	}

	order.Status = entities.OrderStatusShipped

	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return nil, err
	}

	_ = uc.publisher.PublishOrderEvent(ctx, eventbus.OrderEventPayload{
		EventType: eventbus.EventOrderShipped,
		OrderID:   order.ID,
		BuyerID:   order.BuyerID,
		SellerID:  order.SellerID,
		AssetID:   order.AssetID,
	})

	out := response.OrderFromEntity(order)
	return &out, nil
}
