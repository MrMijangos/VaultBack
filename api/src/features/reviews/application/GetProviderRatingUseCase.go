package application

import (
	"context"

	"vault/src/features/reviews/domain/dto/response"
	"vault/src/features/reviews/domain/repositories"
)

type GetProviderRatingUseCase struct {
	repo repositories.ReviewRepository
}

func NewGetProviderRatingUseCase(repo repositories.ReviewRepository) *GetProviderRatingUseCase {
	return &GetProviderRatingUseCase{repo: repo}
}

func (uc *GetProviderRatingUseCase) Execute(ctx context.Context, providerID string) (response.ProviderRatingResponse, error) {
	rating, total, err := uc.repo.GetProviderRating(ctx, providerID)
	if err != nil {
		return response.ProviderRatingResponse{}, err
	}
	return response.ProviderRatingResponse{
		ProviderID:   providerID,
		Rating:       rating,
		TotalReviews: total,
	}, nil
}
