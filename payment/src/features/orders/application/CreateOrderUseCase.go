package application

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"

	"vault-payment/src/core/eventbus"
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
	assetPrices        repositories.AssetPriceProvider
	stripeClient       stripeclient.Client
	publisher          eventbus.Publisher
}

func NewCreateOrderUseCase(
	orderRepo repositories.OrderRepository,
	commissionProvider repositories.SellerCommissionProvider,
	accountProvider repositories.SellerAccountProvider,
	assetPrices repositories.AssetPriceProvider,
	stripeClient stripeclient.Client,
	publisher eventbus.Publisher,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepo:          orderRepo,
		commissionProvider: commissionProvider,
		accountProvider:    accountProvider,
		assetPrices:        assetPrices,
		stripeClient:       stripeClient,
		publisher:          publisher,
	}
}

func (uc *CreateOrderUseCase) Execute(ctx context.Context, buyerID string, req request.CreateOrderRequest) (*response.OrderResponse, error) {
	if req.SellerID == "" || req.AssetID == "" || req.BuyerEmail == "" || req.PaymentMethodID == "" {
		return nil, ErrInvalidRequest
	}
	// AssetID viaja sin validar hasta HTTPAssetPriceProvider, que lo mete
	// directo en la URL de la consulta a api/ -- se exige formato UUID acá
	// para que nadie pueda mandar algo como "../otra-ruta" y desviar esa
	// petición interna a un endpoint distinto.
	if _, err := uuid.Parse(req.AssetID); err != nil {
		return nil, ErrInvalidRequest
	}
	if buyerID == req.SellerID {
		return nil, ErrSelfPurchase
	}

	// El monto a cobrar SIEMPRE se recalcula del lado del servidor a partir
	// del precio real del activo -- antes se cobraba tal cual el
	// amount_cents que mandara el cliente, así que cualquiera podía
	// interceptar la petición y comprar un activo por el precio que
	// quisiera.
	priceCents, isForSale, err := uc.assetPrices.GetSalePriceCents(ctx, req.AssetID)
	if err != nil {
		return nil, err
	}
	if !isForSale || priceCents <= 0 {
		return nil, ErrAssetNotForSale
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

	paymentIntentID, err := uc.stripeClient.ChargeBuyer(ctx, customerID, attachedPaymentMethodID, priceCents, currency)
	if err != nil {
		return nil, fmt.Errorf("no se pudo cobrar al comprador: %w", err)
	}

	commissionCents := int64(math.Round(float64(priceCents) * rate))
	sellerAmountCents := priceCents - commissionCents

	order := &entities.Order{
		ID:                    uuid.NewString(),
		BuyerID:               buyerID,
		SellerID:              req.SellerID,
		AssetID:               req.AssetID,
		AmountCents:           priceCents,
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

	_ = uc.publisher.PublishOrderEvent(ctx, eventbus.OrderEventPayload{
		EventType: eventbus.EventOrderCreated,
		OrderID:   order.ID,
		BuyerID:   order.BuyerID,
		SellerID:  order.SellerID,
		AssetID:   order.AssetID,
	})

	out := response.OrderFromEntity(order)
	return &out, nil
}
