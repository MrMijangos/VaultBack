package repositories

import (
	"context"

	"vault-payment/src/features/subscriptions/domain/entities"
)

type PlanRepository interface {
	List(ctx context.Context) ([]*entities.Plan, error)
	GetByID(ctx context.Context, id string) (*entities.Plan, error)
}
