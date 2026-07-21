package entities

import "time"

const (
	SubscriptionStatusActive   = "active"
	SubscriptionStatusCanceled = "canceled"
	SubscriptionStatusFailed   = "failed"
)

// Subscription es la suscripción vigente (o pasada) de un usuario. Se
// mantiene en memoria por ahora -- BuildInMemorySubscriptionRepository --
// hasta que exista la tabla "subscriptions" en Supabase.
type Subscription struct {
	ID                   string
	UserID               string
	PlanID               string
	Status               string
	StripeCustomerID     string
	StripeSubscriptionID string
	CurrentPeriodStart   time.Time
	CurrentPeriodEnd     time.Time
	CanceledAt           *time.Time
	CreatedAt            time.Time
}

func (s *Subscription) IsActive() bool {
	return s.Status == SubscriptionStatusActive
}
