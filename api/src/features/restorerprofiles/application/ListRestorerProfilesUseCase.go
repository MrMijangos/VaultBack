package application

import (
	"context"

	"vault/src/features/restorerprofiles/domain/dto/response"
	"vault/src/features/restorerprofiles/domain/repositories"
)

type ListRestorerProfilesUseCase struct {
	repo       repositories.RestorerProfileRepository
	ratingRepo repositories.ProviderRatingProvider
}

func NewListRestorerProfilesUseCase(repo repositories.RestorerProfileRepository, ratingRepo repositories.ProviderRatingProvider) *ListRestorerProfilesUseCase {
	return &ListRestorerProfilesUseCase{repo: repo, ratingRepo: ratingRepo}
}

func (uc *ListRestorerProfilesUseCase) Execute(ctx context.Context) ([]response.RestorerProfileResponse, error) {
	list, err := uc.repo.ListWithServices(ctx)
	if err != nil {
		return nil, err
	}

	out := response.FromEntities(list)
	for i, p := range list {
		rating, total, err := uc.ratingRepo.GetProviderRating(ctx, p.UserID)
		if err != nil {
			return nil, err
		}
		out[i].Rating = rating
		out[i].ReviewsCount = total
	}

	return out, nil
}
