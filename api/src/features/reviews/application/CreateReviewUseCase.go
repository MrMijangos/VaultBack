package application

import (
	"context"

	"vault/src/features/reviews/domain/dto/request"
	"vault/src/features/reviews/domain/dto/response"
	"vault/src/features/reviews/domain/entities"
	"vault/src/features/reviews/domain/repositories"
)

type CreateReviewUseCase struct {
	repo repositories.ReviewRepository
}

func NewCreateReviewUseCase(repo repositories.ReviewRepository) *CreateReviewUseCase {
	return &CreateReviewUseCase{repo: repo}
}

func (uc *CreateReviewUseCase) Execute(ctx context.Context, userID string, req request.CreateReviewRequest) (response.ReviewResponse, error) {
	if err := req.Validate(); err != nil {
		return response.ReviewResponse{}, err
	}

	created, err := uc.repo.Create(ctx, entities.Review{
		UserID:     userID,
		ProviderID: req.ProviderID,
		Content:    req.Content,
	})
	if err != nil {
		return response.ReviewResponse{}, err
	}

	return response.FromEntity(created), nil
}
