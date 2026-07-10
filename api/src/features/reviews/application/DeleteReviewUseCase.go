package application

import (
	"context"

	"vault/src/features/reviews/domain/repositories"
)

type DeleteReviewUseCase struct {
	repo repositories.ReviewRepository
}

func NewDeleteReviewUseCase(repo repositories.ReviewRepository) *DeleteReviewUseCase {
	return &DeleteReviewUseCase{repo: repo}
}

func (uc *DeleteReviewUseCase) Execute(ctx context.Context, id string, userID string) error {
	return uc.repo.Delete(ctx, id, userID)
}
