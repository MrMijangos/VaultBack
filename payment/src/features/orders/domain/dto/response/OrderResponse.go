package response

import (
	"time"

	"vault-payment/src/features/orders/domain/entities"
)

type OrderResponse struct {
	ID                string     `json:"id"`
	SellerID          string     `json:"seller_id"`
	AssetID           string     `json:"asset_id"`
	AmountCents       int64      `json:"amount_cents"`
	CommissionCents   int64      `json:"commission_cents"`
	SellerAmountCents int64      `json:"seller_amount_cents"`
	Currency          string     `json:"currency"`
	Status            string     `json:"status"`
	CreatedAt         time.Time  `json:"created_at"`
	ConfirmedAt       *time.Time `json:"confirmed_at,omitempty"`
}

func OrderFromEntity(o *entities.Order) OrderResponse {
	return OrderResponse{
		ID:                o.ID,
		SellerID:          o.SellerID,
		AssetID:           o.AssetID,
		AmountCents:       o.AmountCents,
		CommissionCents:   o.CommissionCents,
		SellerAmountCents: o.SellerAmountCents,
		Currency:          o.Currency,
		Status:            o.Status,
		CreatedAt:         o.CreatedAt,
		ConfirmedAt:       o.ConfirmedAt,
	}
}
