package application

import (
	"context"

	"vault/src/features/recommendations/domain/dto/response"
	"vault/src/features/recommendations/domain/repositories"
)

type GetRecommendationsUseCase struct {
	provider repositories.RecommendationProvider
}

func NewGetRecommendationsUseCase(provider repositories.RecommendationProvider) *GetRecommendationsUseCase {
	return &GetRecommendationsUseCase{provider: provider}
}

func (uc *GetRecommendationsUseCase) Execute(ctx context.Context, userID string, topN int) (response.RecommendationsResponse, error) {
	return uc.provider.GetRecommendedItems(ctx, userID, topN)
}
