package application

import (
	"context"

	"vault/src/features/restorerprofiles/domain/dto/response"
	"vault/src/features/restorerprofiles/domain/repositories"
)

type GetRestorerProfileUseCase struct {
	repo       repositories.RestorerProfileRepository
	ratingRepo repositories.ProviderRatingProvider
}

func NewGetRestorerProfileUseCase(repo repositories.RestorerProfileRepository, ratingRepo repositories.ProviderRatingProvider) *GetRestorerProfileUseCase {
	return &GetRestorerProfileUseCase{repo: repo, ratingRepo: ratingRepo}
}

func (uc *GetRestorerProfileUseCase) Execute(ctx context.Context, userID string) (response.RestorerProfileResponse, error) {
	profile, err := uc.repo.FindByUserID(ctx, userID)
	if err != nil {
		return response.RestorerProfileResponse{}, err
	}

	rating, total, err := uc.ratingRepo.GetProviderRating(ctx, userID)
	if err != nil {
		return response.RestorerProfileResponse{}, err
	}
	profile.Rating = rating
	profile.TotalReviews = total

	return response.FromEntity(profile), nil
}
