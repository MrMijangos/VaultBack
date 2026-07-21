package response

import (
	"time"

	"vault-payment/src/features/subscriptions/domain/entities"
)

type SubscriptionResponse struct {
	ID                 string     `json:"id"`
	PlanID             string     `json:"plan_id"`
	Status             string     `json:"status"`
	CurrentPeriodStart time.Time  `json:"current_period_start"`
	CurrentPeriodEnd   time.Time  `json:"current_period_end"`
	CanceledAt         *time.Time `json:"canceled_at,omitempty"`
}

func SubscriptionFromEntity(s *entities.Subscription) SubscriptionResponse {
	return SubscriptionResponse{
		ID:                 s.ID,
		PlanID:             s.PlanID,
		Status:             s.Status,
		CurrentPeriodStart: s.CurrentPeriodStart,
		CurrentPeriodEnd:   s.CurrentPeriodEnd,
		CanceledAt:         s.CanceledAt,
	}
}
