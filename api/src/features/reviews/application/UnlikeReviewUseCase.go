package application

import (
	"context"

	"vault/src/features/reviews/domain/repositories"
)

type UnlikeReviewUseCase struct {
	repo repositories.ReviewRepository
}

func NewUnlikeReviewUseCase(repo repositories.ReviewRepository) *UnlikeReviewUseCase {
	return &UnlikeReviewUseCase{repo: repo}
}

func (uc *UnlikeReviewUseCase) Execute(ctx context.Context, reviewID string, userID string) error {
	return uc.repo.Unlike(ctx, reviewID, userID)
}
