package application

import (
	"context"

	"vault/src/features/businesses/domain/dto/response"
	"vault/src/features/businesses/domain/repositories"
)

type GetBusinessByIdUseCase struct {
	repo       repositories.BusinessRepository
	ratingRepo repositories.RatingProvider
}

func NewGetBusinessByIdUseCase(repo repositories.BusinessRepository, ratingRepo repositories.RatingProvider) *GetBusinessByIdUseCase {
	return &GetBusinessByIdUseCase{repo: repo, ratingRepo: ratingRepo}
}

func (uc *GetBusinessByIdUseCase) Execute(ctx context.Context, id string) (response.BusinessResponse, error) {
	b, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return response.BusinessResponse{}, err
	}
	out := response.FromEntity(b)
	rating, total, err := uc.ratingRepo.GetProviderRating(ctx, b.UserID)
	if err != nil {
		return response.BusinessResponse{}, err
	}
	out.Rating = rating
	out.TotalReviews = total
	return out, nil
}
