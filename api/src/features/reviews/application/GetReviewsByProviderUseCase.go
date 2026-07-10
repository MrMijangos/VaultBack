package application

import (
	"context"

	"vault/src/features/reviews/domain/dto/response"
	"vault/src/features/reviews/domain/repositories"
)

type GetReviewsByProviderUseCase struct {
	repo repositories.ReviewRepository
}

func NewGetReviewsByProviderUseCase(repo repositories.ReviewRepository) *GetReviewsByProviderUseCase {
	return &GetReviewsByProviderUseCase{repo: repo}
}

func (uc *GetReviewsByProviderUseCase) Execute(ctx context.Context, providerID string) ([]response.ReviewResponse, error) {
	list, err := uc.repo.FindByProviderID(ctx, providerID)
	if err != nil {
		return nil, err
	}
	return response.FromEntities(list), nil
}
