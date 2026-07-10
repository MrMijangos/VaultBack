package application

import (
	"context"

	"vault/src/features/reviews/domain/repositories"
)

type LikeReviewUseCase struct {
	repo repositories.ReviewRepository
}

func NewLikeReviewUseCase(repo repositories.ReviewRepository) *LikeReviewUseCase {
	return &LikeReviewUseCase{repo: repo}
}

func (uc *LikeReviewUseCase) Execute(ctx context.Context, reviewID string, userID string) error {
	return uc.repo.Like(ctx, reviewID, userID)
}
