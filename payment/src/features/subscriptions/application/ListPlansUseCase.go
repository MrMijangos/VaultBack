package application

import (
	"context"

	"vault-payment/src/features/subscriptions/domain/dto/response"
	"vault-payment/src/features/subscriptions/domain/repositories"
)

type ListPlansUseCase struct {
	planRepo repositories.PlanRepository
}

func NewListPlansUseCase(planRepo repositories.PlanRepository) *ListPlansUseCase {
	return &ListPlansUseCase{planRepo: planRepo}
}

func (uc *ListPlansUseCase) Execute(ctx context.Context) ([]response.PlanResponse, error) {
	plans, err := uc.planRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]response.PlanResponse, 0, len(plans))
	for _, p := range plans {
		out = append(out, response.PlanFromEntity(p))
	}
	return out, nil
}
