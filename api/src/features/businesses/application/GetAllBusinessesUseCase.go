package application

import (
	"context"

	"vault/src/features/businesses/domain/dto/response"
	"vault/src/features/businesses/domain/repositories"
)

type GetAllBusinessesUseCase struct {
	repo       repositories.BusinessRepository
	ratingRepo repositories.RatingProvider
}

func NewGetAllBusinessesUseCase(repo repositories.BusinessRepository, ratingRepo repositories.RatingProvider) *GetAllBusinessesUseCase {
	return &GetAllBusinessesUseCase{repo: repo, ratingRepo: ratingRepo}
}

func (uc *GetAllBusinessesUseCase) Execute(ctx context.Context) ([]response.BusinessResponse, error) {
	list, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	out := response.FromEntities(list)
	for i, b := range list {
		rating, total, err := uc.ratingRepo.GetProviderRating(ctx, b.UserID)
		if err != nil {
			return nil, err
		}
		out[i].Rating = rating
		out[i].TotalReviews = total
	}
	return out, nil
}
