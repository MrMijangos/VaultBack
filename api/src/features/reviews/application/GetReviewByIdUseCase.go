package application

import (
	"context"

	"vault/src/features/reviews/domain/dto/response"
	"vault/src/features/reviews/domain/repositories"
)

type GetReviewByIdUseCase struct {
	repo repositories.ReviewRepository
}

func NewGetReviewByIdUseCase(repo repositories.ReviewRepository) *GetReviewByIdUseCase {
	return &GetReviewByIdUseCase{repo: repo}
}

func (uc *GetReviewByIdUseCase) Execute(ctx context.Context, id string) (response.ReviewResponse, error) {
	rv, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return response.ReviewResponse{}, err
	}
	return response.FromEntity(rv), nil
}
